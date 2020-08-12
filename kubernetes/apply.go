package kubernetes

import (
	"context"
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
	Execute(scope string, dryRunOn model.DryRunOn) error
	Wait(scope string, wu model.WaitUntil) (relevantDuration time.Duration, err error)
	Rollback(scope string)
	String() string
}

func NewApplyObject(
	project *model.Project,
	source string,
	object *unstructured.Unstructured,
	client dynamic.Interface,
	runtime Runtime,
) (*ApplyObject, error) {
	objectResource, err := GetObjectResource(object, client, project.Scheme)
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
		object:  objectResource,
		runtime: runtime,
	}, nil
}

type ApplyObject struct {
	log               log.Logger
	KeepAliveInterval time.Duration

	project  *model.Project
	object   ObjectResource
	original *ObjectResource

	applied *unstructured.Unstructured
	runtime Runtime
}

func (instance ApplyObject) String() string {
	return instance.object.String()
}

func (instance *ApplyObject) resolveDryRunOn(in model.DryRunOn) (model.DryRunOn, error) {
	if in == model.DryRunNowhere {
		return model.DryRunNowhere, nil
	}
	ofObject, err := instance.project.Annotations.GetDryRunOnFor(instance.object.Object, in)
	if err != nil {
		return "", err
	}
	return ResolveDryRun(ofObject, instance.object.GroupVersionKind, instance.object.Client, instance.runtime)
}

func (instance *ApplyObject) Execute(scope string, dryRunOn model.DryRunOn) (err error) {
	if dryRunOn, err = instance.resolveDryRunOn(dryRunOn); err != nil {
		return err
	}
	applyOn, err := instance.project.Annotations.GetApplyOnFor(instance.object.Object)
	if err != nil {
		return err
	}
	stage, err := instance.project.Annotations.GetStageFor(instance.object.Object)
	if err != nil {
		return err
	}
	l := instance.log.
		WithField("scope", scope).
		WithField("stage", stage).
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

		return instance.create(scope, dryRunOn)
	} else if err != nil {
		return err
	} else {
		if !applyOn.OnUpdate() {
			l.
				WithField("status", "skipped").
				Debug("%v does exist but should not be updated - skipping.", instance.object)
			return nil
		}

		originalResource, err := GetObjectResource(original, instance.object.Client, instance.project.Scheme)
		if err != nil {
			return err
		}
		instance.original = &originalResource
		l.
			WithField("status", "success").
			WithDeepFieldOn("response", original, l.IsDebugEnabled).
			Debug("%v does exist - it will be updated.", instance.object)

		if err := transformation.Default.TransformForUpdate(instance.project, *original, instance.object.Object); err != nil {
			return err
		}

		return instance.update(scope, *original, dryRunOn)
	}
}

func (instance *ApplyObject) Wait(scope string, global model.WaitUntil) (relevantDuration time.Duration, err error) {
	wu := global
	wuf := wu.AsLazyFormatter("{{with .Timeout}}for {{.}} {{end}}")
	skip := false
	start := time.Now()
	l := instance.log.
		WithField("scope", scope).
		WithField("action", "wait")

	ctx, finished := context.WithCancel(context.Background())

	defer func() {
		if dErr := instance.deleteIfNeeded(scope, wu); dErr != nil {
			if err != nil {
				err = fmt.Errorf("%w - and - %v", err, dErr)
			} else {
				err = dErr
			}
		}
	}()
	defer func() {
		finished()
		ld := l.WithField("duration", time.Now().Sub(start))
		if err != nil {
			ldd := ld.
				WithError(err).
				WithField("status", "failed")
			if ldd.IsDebugEnabled() {
				ldd.Error("Wait %vuntil %v is ready... FAILED!", wuf, instance.object)
			} else {
				ldd.Error("%v is not ready.", instance.object)
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

	if lc := owu.LogConsumer; lc != nil {
		if provider := LogProviderFor(instance.runtime, instance.applied, owu.LogSourceContainerName); provider != nil {
			go func() {
				if err := PrintLogs(ctx, provider, lc.OpenForWrite); err != nil {
					l.WithError(err).WithField("consumer", lc).Error("cannot consume logs")
				}
			}()
		}
	}

	if instance.applied == nil {
		return
	}
	generation := instance.getGenerationOf(instance.applied)
	if generation == nil {
		return 0, fmt.Errorf("cannot retrieve generation of object to be applied")
	}

	resource, rErr := GetObjectResource(instance.applied, instance.object.Client, instance.project.Scheme)
	if rErr != nil {
		return 0, rErr
	}
	for {
		cWu := wu
		if to := wu.Timeout; to != nil {
			timeout := *to - time.Now().Sub(start)
			if timeout <= 0 {
				return 0, common.NewTimeoutError("%v was not ready after %v", resource, *to)
			} else if instance.KeepAliveInterval > 0 && timeout > instance.KeepAliveInterval {
				timeout = instance.KeepAliveInterval
			}
			cWu.Timeout = &timeout
		}
		if done, wErr := instance.watchRun(resource, *generation, cWu, l); wErr != nil || done {
			if owu.Stage == model.WaitUntilStageDefault {
				relevantDuration = time.Now().Sub(start)
			}
			return relevantDuration, wErr
		}
		duration := time.Now().Sub(start)
		l.
			WithField("duration", duration).
			WithField("status", "continue").
			Info("%v is still not ready after %v. Continue waiting...", resource, duration)
	}
}

func (instance *ApplyObject) watchRun(resource ObjectResource, generation int64, wu model.WaitUntil, l log.Logger) (done bool, err error) {
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
	if timeout := wu.Timeout; timeout != nil && *timeout > 0 {
		start := time.Now()
		for {
			cTimeout := *timeout - time.Now().Sub(start)
			if cTimeout <= 0 {
				return false, nil
			}
			select {
			case event := <-w.ResultChan():
				if done, oErr := instance.onWatchEvent(event, l, generation, wu.Stage); oErr != nil || done {
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
		if done, oErr := instance.onWatchEvent(event, l, generation, wu.Stage); oErr != nil || done {
			return done, oErr
		}
	}
}

func (instance *ApplyObject) onWatchEvent(event watch.Event, l log.Logger, generation int64, wus model.WaitUntilStage) (done bool, err error) {
	objectInfo, _ := GetObjectInfo(event.Object, instance.project.Scheme)
	l = l.WithDeepFieldOn("event", event, log.IsTraceEnabled)
	l.Trace("Received event %v on %v.", event.Type, objectInfo)

	if !instance.matchesReferenceOfObjectToApplyAndGeneration(event.Object, generation) {
		l.Trace("Received event %v on %v which does not match %v and will be ignored.", event.Type, objectInfo, instance.object)
		return false, nil
	}

	switch wus {
	case model.WaitUntilStageApplied:
		return instance.onWatchEventForApplied(event, objectInfo, l)
	case model.WaitUntilStageExecuted:
		return instance.onWatchEventForExecuted(event, objectInfo, l)
	default:
		return true, fmt.Errorf("at this position waitUntil.stage of '%v' is not expected", wus)
	}
}

func (instance *ApplyObject) onWatchEventForApplied(event watch.Event, objectInfo ObjectInfo, l log.Logger) (done bool, err error) {
	if ready := IsReady(event.Object); ready == nil {
		l.Debug("Received event %v on %v does not support ready check and will be assumed as ready now.", event.Type, objectInfo)
		return true, nil
	} else if *ready {
		l.Debug("Received event %v on %v which passes the ready check.", event.Type, objectInfo)
		return true, nil
	}
	l.Debug("Received event %v on %v which does not pass the ready check. Continue wait...", event.Type, objectInfo)
	return false, nil
}

func (instance *ApplyObject) onWatchEventForExecuted(event watch.Event, objectInfo ObjectInfo, l log.Logger) (done bool, err error) {
	unknownFail := func() (done bool, err error) {
		return true, fmt.Errorf("don't know how to watch for executed stage of object")
	}
	if state := StateOf(event.Object); state == nil {
		return unknownFail()
	} else if state.IsActive() {
		l.Debug("Received event %v on %v which does indicate that the object is still active. Continue wait...", event.Type, objectInfo)
		return false, nil
	} else if *state == StateSucceeded {
		l.Debug("Received event %v on %v which passes the ready check.", event.Type, objectInfo)
		return true, nil
	} else if *state == StateFailed {
		return true, fmt.Errorf("execution failed")
	}
	return unknownFail()
}

func (instance *ApplyObject) create(scope string, dry model.DryRunOn) (err error) {
	start := time.Now()
	l := instance.log.
		WithField("scope", scope).
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

	target, cErr := instance.object.CloneForCreate(instance.project)
	if cErr != nil {
		return cErr
	}

	opts := metav1.CreateOptions{}
	if dry == model.DryRunOnServer {
		opts.DryRun = []string{metav1.DryRunAll}
	}

	if dry != model.DryRunOnClient {
		if instance.applied, err = target.Create(&opts); err != nil {
			instance.applied = nil
			return
		}
	}
	return
}

func (instance *ApplyObject) update(scope string, original unstructured.Unstructured, dry model.DryRunOn) (err error) {
	start := time.Now()
	l := instance.log.
		WithField("scope", scope).
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

	target, cErr := instance.object.CloneForUpdate(instance.project, original)
	if cErr != nil {
		return cErr
	}

	opts := metav1.UpdateOptions{}
	if dry == model.DryRunOnServer {
		opts.DryRun = []string{"All"}
	}
	if dry != model.DryRunOnClient {
		if instance.applied, err = target.Update(&opts); err != nil {
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

func (instance *ApplyObject) Delete(scope string) (err error) {
	start := time.Now()
	l := instance.log.
		WithField("scope", scope).
		WithField("action", "delete")

	defer func() {
		ld := l.WithField("duration", time.Now().Sub(start))
		if err != nil {
			ldd := ld.
				WithError(err).
				WithField("status", "failed")
			if ldd.IsDebugEnabled() {
				ldd.Error("Deleting %v... FAILED!", instance.object)
			} else {
				ldd.Error("Was not able to delete %v.", instance.object)
			}
		} else {
			ldd := ld.WithField("status", "success")
			if ldd.IsDebugEnabled() {
				ldd.Info("Deleting %v... DONE!", instance.object)
			} else {
				ldd.Info("%v deleted.", instance.object)
			}
		}
	}()

	l.Debug("Deleting %v...", instance.object)

	dp := metav1.DeletePropagationForeground
	if err := instance.object.Delete(&metav1.DeleteOptions{
		PropagationPolicy: &dp,
	}); err != nil {
		return fmt.Errorf("cannot delete resource: %w", err)
	}

	return nil
}

func (instance *ApplyObject) deleteIfNeeded(scope string, cu model.WaitUntil) error {
	if cu.Stage != model.WaitUntilStageExecuted {
		return nil
	}

	cleanupOn, err := instance.project.Annotations.GetCleanupOn(instance.object.Object)
	if err != nil {
		return fmt.Errorf("cannot delete resource: %w", err)
	}

	if !cleanupOn.OnExecuted() {
		return nil
	}

	return instance.Delete(scope)
}

func (instance *ApplyObject) Rollback(scope string) {
	if instance.applied == nil {
		return
	}
	var err error
	start := time.Now()
	l := instance.log.WithField("action", "rollback").WithField("scope", scope)
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

func (instance ApplySet) Execute(scope string, dryRunOn model.DryRunOn) (err error) {
	defer func() {
		if err != nil && dryRunOn == model.DryRunNowhere {
			instance.Rollback(scope)
		}
	}()
	for _, child := range instance {
		if err = child.Execute(scope, dryRunOn); err != nil {
			err = fmt.Errorf("cannot apply %v: %w", child, err)
			return
		}
	}
	return
}

func (instance ApplySet) Rollback(scope string) {
	for _, action := range instance {
		action.Rollback(scope)
	}
}

func (instance ApplySet) Wait(scope string, wu model.WaitUntil) (relevantDuration time.Duration, err error) {
	defer func() {
		if err != nil {
			instance.Rollback(scope)
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
		if cRelevantDuration, cErr := child.Wait(scope, cWu); cErr != nil {
			return 0, fmt.Errorf("cannot wait for %v: %w", child, cErr)
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
	if instance == nil || *instance == nil {
		*instance = StagedApplySet{}
	}
	set := (*instance)[stage]
	set.Add(apply)
	(*instance)[stage] = set
}

func (instance StagedApplySet) Execute(scope string, dry model.DryRunOn, wu *model.WaitUntil, rollbackIfNeeded bool) (relevantDuration time.Duration, err error) {
	defer func() {
		if err != nil && rollbackIfNeeded {
			for _, action := range instance {
				action.Rollback(scope)
			}
		}
	}()
	for stage := range instance {
		cWu := wu
		if cWu != nil && cWu.Timeout != nil {
			if relevantDuration > *cWu.Timeout {
				return 0, common.NewTimeoutError("timeout of %v reached - no more time to continue with left resources", *cWu.Timeout)
			}
			cTimeout := *cWu.Timeout - relevantDuration
			tcWu := wu.CopyWithTimeout(&cTimeout)
			cWu = &tcWu
		}
		if eRelevantDuration, eErr := instance.ExecuteStage(scope, stage, dry, cWu); eErr != nil {
			return 0, eErr
		} else {
			relevantDuration += eRelevantDuration
		}
	}
	return
}

func (instance StagedApplySet) ExecuteStage(scope string, stage model.Stage, dryRunOn model.DryRunOn, wu *model.WaitUntil) (relevantDuration time.Duration, err error) {
	set := instance[stage]
	start := time.Now()
	l := log.WithField("stage", stage).
		WithField("scope", scope)

	l.Info("Entering %s/%v...", scope, stage)
	defer func() {
		l = l.WithField("duration", time.Now().Sub(start))
		if err != nil {
			set.Rollback(scope)
			l.WithError(err).Error("Entering %s/%v... FAILED!", scope, stage)
		} else {
			l.Debug("Entering %s/%v... SUCCESS!", scope, stage)
		}
	}()
	if eErr := set.Execute(scope, dryRunOn); eErr != nil {
		return 0, eErr
	}
	if wu != nil {
		relevantDuration, err = set.Wait(scope, *wu)
	}
	return
}
