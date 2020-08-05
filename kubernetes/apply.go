package kubernetes

import (
	"fmt"
	"github.com/echocat/kubor/common"
	"github.com/echocat/kubor/kubernetes/transformation"
	"github.com/echocat/kubor/log"
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"reflect"
	"time"
)

type Apply interface {
	Execute(DryRunOn) error
	Wait(wu model.WaitUntil) (relevantDuration time.Duration, err error)
	Rollback()
	String() string
}

func NewApplyObject(
	project *model.Project,
	source string,
	object *unstructured.Unstructured,
	client dynamic.Interface,
	runtime Runtime,
	objectValidator ObjectValidator,
) (*ApplyObject, error) {
	objectResource, err := GetObjectResource(object, client, objectValidator)
	if err != nil {
		return nil, err
	}

	stage, err := project.Annotations.GetStageFor(object)
	if err != nil {
		return nil, err
	}

	return &ApplyObject{
		project: project,
		log: log.
			WithField("source", source).
			WithField("object", objectResource).
			WithField("stage", stage),
		object:          objectResource,
		runtime:         runtime,
		objectValidator: objectValidator,
	}, nil
}

type ApplyObject struct {
	log               log.Logger
	KeepAliveInterval time.Duration

	project  *model.Project
	object   ObjectResource
	original *ObjectResource

	applied         *unstructured.Unstructured
	objectValidator ObjectValidator
	runtime         Runtime
}

func (instance ApplyObject) String() string {
	return instance.object.String()
}

func (instance *ApplyObject) resolveDryRunOn(dry DryRunOn) (DryRunOn, error) {
	return dry.Resolve(instance.object.Kind, instance.object.Client, instance.runtime)
}

func (instance *ApplyObject) Execute(dry DryRunOn) (err error) {
	if dry, err = instance.resolveDryRunOn(dry); err != nil {
		return err
	}
	applyOn, err := instance.project.Annotations.GetApplyOnFor(instance.object.Object)
	if err != nil {
		return err
	}
	l := instance.log.
		WithField("action", "checkExistence")
	original, err := instance.object.Get(nil)
	if errors.IsNotFound(err) {
		if !applyOn.OnCreate() {
			l.
				WithField("status", "skipped").
				Debug("%v does not exist but should not be created - skipping.", instance.object)
			return nil
		}

		l.
			WithField("status", "notFound").
			Debug("%v does not exist - it will be created.", instance.object)
		instance.original = nil

		if err := transformation.TransformForCreate(instance.project, instance.object.Object); err != nil {
			return err
		}

		return instance.create(dry)
	} else if err != nil {
		return err
	} else {
		if !applyOn.OnUpdate() {
			l.
				WithField("status", "skipped").
				Debug("%v does exist but should not be updated - skipping.", instance.object)
			return nil
		}

		originalResource, err := GetObjectResource(original, instance.object.Client, instance.objectValidator)
		if err != nil {
			return err
		}
		instance.original = &originalResource
		l.
			WithField("status", "success").
			WithDeepFieldOn("response", original, l.IsDebugEnabled).
			Debug("%v does exist - it will be updated.", instance.object)

		if err := transformation.TransformForUpdate(instance.project, *original, instance.object.Object); err != nil {
			return err
		}

		return instance.update(dry)
	}
}

func (instance *ApplyObject) Wait(global model.WaitUntil) (relevantDuration time.Duration, err error) {
	wu := global
	wuf := wu.AsLazyFormatter("{{with .Timeout}}for {{.}} {{end}}")
	skip := false
	start := time.Now()
	l := instance.log.
		WithField("action", "wait")
	defer func() {
		ld := l.WithField("duration", time.Now().Sub(start))
		if err != nil {
			ldd := ld.
				WithError(err).
				WithField("status", "failed")
			if ldd.IsDebugEnabled() {
				ldd.Error("Wait %vuntil %v is ready... FAILED!", wuf, instance.object)
			} else {
				ldd.Error("%v was not ready after %v.", instance.object, wuf)
			}
		} else if skip {
			ldd := ld.WithField("status", "skipped")
			if ldd.IsDebugEnabled() {
				ldd.Info("Wait %vuntil %v is ready... SKIPPED!", wuf, instance.object)
			}
		} else {
			ldd := ld.WithField("status", "success")
			if ldd.IsDebugEnabled() {
				ldd.Info("Wait %vuntil %v is ready... DONE!", wuf, instance.object)
			} else {
				ldd.Info("%v is ready.", instance.object)
			}
		}
	}()

	owu, wuErr := instance.project.Annotations.GetWaitUntilFor(instance.object.Object)
	if wuErr != nil {
		return 0, wuErr
	}
	wu = wu.MergeWith(owu)
	wuf.WaitUntil = wu

	if to := wu.Timeout; to != nil {
		l = l.WithField("timeout", *to)
	} else {
		l = l.WithField("timeout", "unlimited")
	}
	l.Debug("Wait %vuntil %v is ready...", wuf, instance.object)

	if wu.Stage == model.WaitUntilStageNever {
		skip = true
		return
	}

	if instance.applied == nil {
		return
	}
	generation := instance.getGenerationOf(instance.applied)
	if generation == nil {
		return 0, fmt.Errorf("cannot retrieve generation of object to be applied")
	}

	resource, rErr := GetObjectResource(instance.applied, instance.object.Client, instance.objectValidator)
	if rErr != nil {
		return 0, rErr
	}
	for {
		var timeout time.Duration
		if to := wu.Timeout; to != nil {
			timeout = *to - time.Now().Sub(start)
			if timeout <= 0 {
				return 0, common.NewTimeoutError("%v was not ready after %v", resource, *to)
			} else if instance.KeepAliveInterval > 0 && timeout > instance.KeepAliveInterval {
				timeout = instance.KeepAliveInterval
			}
		}
		if done, wErr := instance.watchRun(resource, *generation, timeout, l); wErr != nil || done {
			if owu.Stage == model.WaitUntilStageDefault {
				relevantDuration = time.Now().Sub(start)
			}
			return 0, wErr
		}
		duration := time.Now().Sub(start)
		l.
			WithField("duration", duration).
			WithField("status", "continue").
			Info("%v is still not ready after %v. Continue waiting...", resource, duration)
	}
}

func (instance *ApplyObject) watchRun(resource ObjectResource, generation int64, timeout time.Duration, l log.Logger) (done bool, err error) {
	w, wErr := resource.Watch(nil)
	if wErr != nil {
		return false, wErr
	}
	defer w.Stop()
	get, err := resource.Get(nil)
	if err != nil {
		return false, err
	}
	if instance.matchesReferenceOfObjectToApplyAndGenerationAndIsReady(get, generation) {
		return true, nil
	}
	if timeout > 0 {
		start := time.Now()
		for {
			cTimeout := timeout - time.Now().Sub(start)
			if cTimeout <= 0 {
				return false, nil
			}
			select {
			case event := <-w.ResultChan():
				if done, oErr := instance.onWatchEvent(event, l, generation); oErr != nil || done {
					return done, oErr
				}
			case <-time.After(cTimeout):
				get, err := resource.Get(nil)
				if err != nil {
					return false, err
				}
				return instance.matchesReferenceOfObjectToApplyAndGenerationAndIsReady(get, generation), nil
			}
		}
	}

	for {
		event := <-w.ResultChan()
		if done, oErr := instance.onWatchEvent(event, l, generation); oErr != nil || done {
			return done, oErr
		}
	}
}

func (instance *ApplyObject) onWatchEvent(event watch.Event, l log.Logger, generation int64) (done bool, err error) {
	eventObjectInfo, _ := GetObjectInfo(event.Object, instance.objectValidator)
	ld := l.WithDeepFieldOn("event", event, log.IsTraceEnabled)
	ld.WithField("event", event).
		Trace("Received event %v on %v.", event.Type, eventObjectInfo)

	if !instance.matchesReferenceOfObjectToApplyAndGeneration(event.Object, generation) {
		ld.Trace("Received event %v on %v which does not match %v and will be ignored.", event.Type, eventObjectInfo, instance.object)
	} else if ready := IsReady(event.Object); ready == nil {
		ld.Debug("Received event %v on %v does not support ready check and will be assumed as ready now.", event.Type, eventObjectInfo)
		return true, nil
	} else if *ready {
		ld.Debug("Received event %v on %v which passes the ready check.", event.Type, eventObjectInfo)
		return true, nil
	} else {
		ld.Debug("Received event %v on %v which does not pass the ready check. Continue wait...", event.Type, eventObjectInfo)
	}
	return false, nil
}

func (instance *ApplyObject) create(dry DryRunOn) (err error) {
	start := time.Now()
	l := instance.log.
		WithField("action", "create").
		WithField("dryRunOn", dry)
	defer func() {
		ld := l.
			WithField("duration", time.Now().Sub(start)).
			WithDeepFieldOn("response", instance.applied, l.IsTraceEnabled)
		if err != nil {
			ldd := ld.
				WithError(err).
				WithField("status", "failed")
			if ldd.IsDebugEnabled() {
				ldd.Error("Create %v... FAILED!", instance.object)
			} else {
				ldd.Error("Could not create %v.", instance.object)
			}
		} else {
			ldd := ld.
				WithField("status", "success")
			if ldd.IsDebugEnabled() {
				ldd.Info("Create %v... SUCCESS!", instance.object)
			} else {
				ldd.Info("%v created.", instance.object)
			}
		}
	}()
	l.Debug("Create %v...", instance.object)
	opts := metav1.CreateOptions{}
	if dry == ServerDryRun {
		opts.DryRun = []string{metav1.DryRunAll}
	}
	if dry != ClientDryRun {
		if instance.applied, err = instance.object.Create(&opts); err != nil {
			instance.applied = nil
			return
		}
	}
	return
}

func (instance *ApplyObject) update(dry DryRunOn) (err error) {
	start := time.Now()
	l := instance.log.
		WithField("action", "update").
		WithField("dryRunOn", dry)
	defer func() {
		ld := l.
			WithField("duration", time.Now().Sub(start)).
			WithDeepFieldOn("response", instance.applied, l.IsTraceEnabled)
		if err != nil {
			ldd := ld.
				WithError(err).
				WithField("status", "failed")
			if ldd.IsDebugEnabled() {
				ldd.Error("Update %v... FAILED!", instance.object)
			} else {
				ldd.Error("Could not update %v.", instance.object)
			}
		} else {
			ldd := ld.
				WithField("status", "success")
			if ldd.IsDebugEnabled() {
				ldd.Info("Update %v... SUCCESS!", instance.object)
			} else {
				ldd.Info("%v updated.", instance.object)
			}
		}
	}()
	l.Debug("Update %v...", instance.object)
	opts := metav1.UpdateOptions{}
	if dry == ServerDryRun {
		opts.DryRun = []string{"All"}
	}
	if dry != ClientDryRun {
		if instance.applied, err = instance.object.Update(&opts); err != nil {
			instance.applied = nil
			return
		}
	}
	return
}

func (instance *ApplyObject) matchesReferenceOfObjectToApplyAndGenerationAndIsReady(runtimeObject runtime.Object, expectedGeneration int64) bool {
	if !instance.matchesReferenceOfObjectToApplyAndGeneration(runtimeObject, expectedGeneration) {
		return false
	}
	ready := IsReady(runtimeObject)
	return ready == nil || *ready
}

func (instance *ApplyObject) matchesReferenceOfObjectToApplyAndGeneration(runtimeObject runtime.Object, expectedGeneration int64) bool {
	actualGeneration := instance.getGenerationOf(runtimeObject)
	if actualGeneration == nil || *actualGeneration != expectedGeneration {
		return false
	}
	return instance.matchesReferenceOfObjectToApply(runtimeObject)
}

func (instance *ApplyObject) matchesReferenceOfObjectToApply(runtimeObject runtime.Object) bool {
	if metaObject, ok := runtimeObject.(metav1.Object); !ok {
		return false
	} else {
		if !reflect.DeepEqual(
			instance.object.Object.GroupVersionKind(),
			runtimeObject.GetObjectKind().GroupVersionKind(),
		) {
			return false
		}

		if !reflect.DeepEqual(
			instance.object.Object.GetNamespace(),
			metaObject.GetNamespace(),
		) {
			return false
		}

		if !reflect.DeepEqual(
			instance.object.Object.GetName(),
			metaObject.GetName(),
		) {
			return false
		}

		return true
	}
}

func (instance *ApplyObject) getGenerationOf(runtimeObject runtime.Object) *int64 {
	if metaObject, ok := runtimeObject.(metav1.Object); !ok {
		return nil
	} else {
		generation := metaObject.GetGeneration()
		return &generation
	}
}

func (instance *ApplyObject) Rollback() {
	if instance.applied == nil {
		return
	}
	var err error
	start := time.Now()
	l := instance.log.WithField("action", "rollback")
	defer func() {
		instance.applied = nil
		ld := l.
			WithField("duration", time.Now().Sub(start))
		if err != nil {
			ldd := ld.
				WithError(err).
				WithField("status", "failed")
			if ldd.IsDebugEnabled() {
				ldd.Warn("Rollback %v... FAILED!", instance.object)
			} else {
				ldd.Warn("Could not rollback %v.", instance.object)
			}
		} else {
			ldd := ld.
				WithField("status", "success")
			if ldd.IsDebugEnabled() {
				ldd.Info("Rollback %v... SUCCESS!", instance.object)
			} else {
				ldd.Info("%v rolled back.", instance.object)
			}
		}
	}()
	l.Debug("Rollback %v...", instance.object)
	if instance.original == nil {
		err = instance.object.Delete(nil)
	} else {
		_, err = instance.original.Update(nil)
	}
}

type ApplySet []Apply

func (instance *ApplySet) Add(apply Apply) {
	if instance == nil {
		*instance = ApplySet{}
	}
	*instance = append(*instance, apply)
}

func (instance ApplySet) Execute(dry DryRunOn) (err error) {
	defer func() {
		if err != nil && dry == NowhereDryRun {
			instance.Rollback()
		}
	}()
	for _, child := range instance {
		if err = child.Execute(dry); err != nil {
			err = fmt.Errorf("cannot apply %v: %v", child, err)
			return
		}
	}
	return
}

func (instance ApplySet) Rollback() {
	for _, action := range instance {
		action.Rollback()
	}
}

func (instance ApplySet) Wait(wu model.WaitUntil) (relevantDuration time.Duration, err error) {
	defer func() {
		if err != nil {
			instance.Rollback()
		}
	}()
	for _, child := range instance {
		cWu := wu
		if to := cWu.Timeout; to != nil {
			if relevantDuration > *to {
				return 0, common.NewTimeoutError("timeout of %v reached - no more time to continue with left resources", *to)
			}
			cTimeout := *to - relevantDuration
			cWu = wu.CopyWithTimeout(&cTimeout)
		}
		if cRelevantDuration, cErr := child.Wait(cWu); cErr != nil {
			return 0, fmt.Errorf("cannot wait for %v: %v", child, err)
		} else {
			relevantDuration += cRelevantDuration
		}
	}
	return
}

func (instance ApplySet) String() string {
	var result string
	for i, child := range instance {
		if i > 0 {
			result += ", "
		}
		result += child.String()
	}
	return "[" + result + "]"
}

type StagedApplySet map[model.Stage]ApplySet

func (instance *StagedApplySet) Add(stage model.Stage, apply Apply) {
	if instance == nil {
		*instance = StagedApplySet{}
	}
	set := (*instance)[stage]
	set.Add(apply)
	(*instance)[stage] = set
}

func (instance StagedApplySet) Execute(wu model.WaitUntil) (relevantDuration time.Duration, err error) {
	defer func() {
		if err != nil {
			instance.Rollback()
		}
	}()
	for stage := range instance {
		cWu := wu
		if to := cWu.Timeout; to != nil {
			if relevantDuration > *to {
				return 0, common.NewTimeoutError("timeout of %v reached - no more time to continue with left resources", *to)
			}
			cTimeout := *to - relevantDuration
			cWu = wu.CopyWithTimeout(&cTimeout)
		}
		if eRelevantDuration, eErr := instance.ExecuteStage(stage, cWu); eErr != nil {
			return 0, eErr
		} else {
			relevantDuration += eRelevantDuration
		}
	}
	return
}

func (instance StagedApplySet) DryRun(dry DryRunOn) error {
	if dry == NowhereDryRun {
		return nil
	}
	for _, child := range instance {
		if err := child.Execute(dry); err != nil {
			return err
		}
	}
	return nil
}

func (instance StagedApplySet) ExecuteStage(stage model.Stage, wu model.WaitUntil) (relevantDuration time.Duration, err error) {
	set := instance[stage]
	defer func() {
		if err != nil {
			set.Rollback()
		}
	}()
	if eErr := set.Execute(NowhereDryRun); eErr != nil {
		return 0, eErr
	}
	return set.Wait(wu)
}

func (instance StagedApplySet) Rollback() {
	for _, action := range instance {
		action.Rollback()
	}
}
