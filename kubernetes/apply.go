package kubernetes

import (
	"fmt"
	"github.com/levertonai/kubor/common"
	"github.com/levertonai/kubor/kubernetes/fixes"
	"github.com/levertonai/kubor/log"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"reflect"
	"time"
)

type Apply interface {
	Execute(DryRunOn) error
	Wait(timeout time.Duration) error
	Rollback()
	String() string
}

func NewApplyObject(source string, object *unstructured.Unstructured, client dynamic.Interface, runtime Runtime) (*ApplyObject, error) {
	objectResource, err := GetObjectResource(object, client)
	if err != nil {
		return nil, err
	}

	return &ApplyObject{
		log: log.
			WithField("source", source).
			WithField("object", objectResource),
		object:  objectResource,
		runtime: runtime,
	}, nil
}

type ApplyObject struct {
	log log.Logger

	object   ObjectResource
	original *ObjectResource

	applied *unstructured.Unstructured
	runtime Runtime
}

func (instance ApplyObject) String() string {
	return instance.object.String()
}

func (instance *ApplyObject) Execute(dry DryRunOn) error {
	if dry == ServerIfPossibleDryRun || dry == ServerDryRun {
		if serverSidePossible, err := HasServerDryRunSupport(instance.object.Kind, instance.object.Client, instance.runtime); err != nil {
			return err
		} else if dry == ServerDryRun {
			if !serverSidePossible {
				return fmt.Errorf("%v does not support server side dry run", instance.object.Kind)
			}
		} else if serverSidePossible {
			dry = ServerDryRun
		} else {
			dry = ClientDryRun
		}
	}
	l := instance.log.
		WithField("action", "checkExistence")
	var err error
	original, err := instance.object.Get(nil)
	if errors.IsNotFound(err) {
		l.
			WithField("status", "notFound").
			Debug("%v does not exist - it will be created.", instance.object)
		instance.original = nil

		if err := fixes.FixForCreate(instance.object.Object); err != nil {
			return err
		}

		return instance.create(dry)
	} else if err != nil {
		return err
	} else {
		originalResource, err := GetObjectResource(original, instance.object.Client)
		if err != nil {
			return err
		}
		instance.original = &originalResource
		l.
			WithField("status", "success").
			WithDeepFieldOn("response", original, l.IsDebugEnabled).
			Debug("%v does exist - it will be updated.", instance.object)

		if err := fixes.FixForUpdate(*original, instance.object.Object); err != nil {
			return err
		}

		return instance.update(dry)
	}
}

func (instance *ApplyObject) Wait(timeout time.Duration) (err error) {
	start := time.Now()
	l := instance.log.
		WithField("action", "wait").
		WithField("timeout", timeout)
	defer func() {
		ld := l.WithField("duration", time.Now().Sub(start))
		if err != nil {
			ldd := ld.
				WithError(err).
				WithField("status", "failed")
			if ldd.IsDebugEnabled() {
				ldd.Error("Wait for %v until %v is ready... FAILED!", timeout, instance.object)
			} else {
				ldd.Error("%v was not ready after %v.", instance.object, timeout)
			}
		} else {
			ldd := ld.
				WithField("status", "success")
			if ldd.IsDebugEnabled() {
				ldd.Info("Wait for %v until %v is ready... DONE!", timeout, instance.object)
			} else {
				ldd.Info("%v is ready.", instance.object)
			}
		}
	}()
	l.Debug("Wait for %v until %v is ready...", timeout, instance.object)

	if instance.applied == nil {
		return
	}
	generation := instance.getGenerationOf(instance.applied)
	if generation == nil {
		return fmt.Errorf("cannot retrieve generation of object to be applied")
	}

	resource, rErr := GetObjectResource(instance.applied, instance.object.Client)
	if rErr != nil {
		err = rErr
		return
	}
	w, wErr := resource.Watch(nil)
	if wErr != nil {
		err = wErr
		return
	}
	get, gErr := resource.Get(nil)
	if gErr != nil {
		err = gErr
		return
	}
	if instance.matchesReferenceOfObjectToApplyAndGenerationAndIsReady(get, *generation) {
		return
	}
	rc := w.ResultChan()
	for afterCh := time.After(timeout); ; {
		select {
		case event := <-rc:
			eventObjectInfo, _ := GetObjectInfo(event.Object)
			ld := log.
				WithDeepFieldOn("event", event, log.IsTraceEnabled)

			if !instance.matchesReferenceOfObjectToApplyAndGeneration(event.Object, *generation) {
				ld.Trace("Received event %v on %v which does not match %v and will be ignored.", event.Type, eventObjectInfo, instance.object)
			} else if ready := IsReady(event.Object); ready == nil {
				ld.Debug("Received event %v on %v does not support ready check and will be assumed as ready now.", event.Type, eventObjectInfo)
				return
			} else if *ready {
				ld.Debug("Received event %v on %v which passes the ready check.", event.Type, eventObjectInfo)
				return
			} else {
				ld.Debug("Received event %v on %v which does not pass the ready check. Continue wait...", event.Type, eventObjectInfo)
			}
		case <-afterCh:
			err = common.NewTimeoutError("%v was not ready after %v", resource, timeout)
			return
		}
	}
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
		if err != nil && !dry.IsEnabled() {
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

func (instance ApplySet) Wait(timeout time.Duration) (err error) {
	defer func() {
		if err != nil {
			instance.Rollback()
		}
	}()
	start := time.Now()
	for _, child := range instance {
		cTimeout := timeout - time.Now().Sub(start)
		if err = child.Wait(cTimeout); err != nil {
			err = fmt.Errorf("cannot wait for %v: %v", child, err)
			return
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
