package kubernetes

import (
	"fmt"
	"github.com/echocat/kubor/log"
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"time"
)

type CleanupTask struct {
	project *model.Project
	keep    gvked
	client  dynamic.Interface
}

func NewCleanupTask(project *model.Project, client dynamic.Interface) (CleanupTask, error) {
	return CleanupTask{
		project: project,
		client:  client,
	}, nil
}

func (instance *CleanupTask) Add(reference model.ObjectReference) {
	instance.keep.add(reference)
}

func (instance *CleanupTask) Execute() error {
	namespaces, err := instance.getNamespaces()
	if err != nil {
		return err
	}

	for _, namespace := range namespaces {
		if err := instance.ExecuteIn(namespace); err != nil {
			return err
		}
	}

	return nil
}

func (instance *CleanupTask) ExecuteIn(namespace model.Namespace) (err error) {
	l := log.WithField("namespace", namespace)

	start := time.Now()

	defer func() {
		l = l.WithField("duration", time.Now().Sub(start))
		if err != nil {
			if l.IsDebugEnabled() {
				l.WithError(err).Debug("Cleanup namespace %v if required... FAILED!", namespace)
			} else {
				l.WithError(err).Error("Cleanup namespace %v failed.", namespace)
			}
		} else {
			if l.IsDebugEnabled() {
				l.Info("Cleanup namespace %v if required... FINISHED!", namespace)
			} else {
				l.Info("Namespace %v is now clean.", namespace)
			}
		}
	}()

	l.Debug("Cleanup namespace %v if required...", namespace)

	for gvk := range instance.project.Claim.GroupVersionKinds {
		if err := instance.executeFor(l, namespace, gvk); err != nil {
			return err
		}
	}

	return nil
}

func (instance *CleanupTask) executeFor(l log.Logger, namespace model.Namespace, gvk model.GroupVersionKind) (err error) {
	l = l.WithField("gvk", gvk)

	start := time.Now()

	defer func() {
		l = l.WithField("duration", time.Now().Sub(start))
		if err != nil {
			l.WithError(err).Trace("Check %v in %v if resources needs to be removed... FAILED!", gvk, namespace)
		} else {
			l.Trace("Check %v in %v if resources needs to be removed... FINISHED!", gvk, namespace)
		}
	}()

	l.Trace("Check %v in %v if resources needs to be removed...", gvk, namespace)

	labelSelector := instance.labelSelector()
	gvr, _ := gvk.GuessToResource()
	resource := instance.client.Resource(gvr).Namespace(namespace.String())
	opts := metav1.ListOptions{
		LabelSelector: labelSelector,
	}
	for {
		list, err := resource.List(opts)
		if err != nil {
			if as, ok := err.(errors.APIStatus); ok && as.Status().Code == 404 {
				return nil
			}
			return fmt.Errorf("cannot collect existing elements of type %v: %w", gvr, err)
		}

		for _, candidate := range list.Items {
			reference, err := GetObjectReference(&candidate, instance.project.Scheme)
			if err != nil {
				l.WithError(err).
					WithField("reference", fmt.Sprintf("%v %v/%v", candidate.GroupVersionKind(), candidate.GetName(), candidate.GetNamespace())).
					Warn("Cannot evaluate object. Skipping it...")
				continue
			}
			if instance.keep.has(reference) {
				l.WithField("reference", reference).
					Trace("Element %v will be kept.", reference)
				continue
			}
			//TODO! Check if we're allowed to delete it.
			if err := instance.delete(resource, reference); err != nil {
				return err
			}
		}

		if v := list.GetContinue(); v != "" {
			opts.Continue = v
		} else {
			return nil
		}
	}
}

func (instance *CleanupTask) delete(resource dynamic.ResourceInterface, reference model.ObjectReference) (err error) {
	start := time.Now()
	l := log.
		WithField("action", "delete")

	defer func() {
		ld := l.WithField("duration", time.Now().Sub(start))
		if err != nil {
			ldd := ld.
				WithError(err).
				WithField("status", "failed")
			if ldd.IsDebugEnabled() {
				ldd.Error("Deleting orphan %v... FAILED!", reference)
			} else {
				ldd.Error("Was not able to delete orphan %v.", reference)
			}
		} else {
			ldd := ld.WithField("status", "success")
			if ldd.IsDebugEnabled() {
				ldd.Info("Deleting orphan %v... DONE!", reference)
			} else {
				ldd.Info("Orphan %v deleted.", reference)
			}
		}
	}()

	l.Debug("Deleting orphan %v...", reference)

	dp := metav1.DeletePropagationForeground
	if err := resource.Delete(reference.Name.String(), &metav1.DeleteOptions{
		PropagationPolicy: &dp,
	}); err != nil {
		return fmt.Errorf("cannot delete resource: %w", err)
	}

	return nil
}

func (instance *CleanupTask) labelSelector() string {
	return fmt.Sprintf("%v=%v,%v=%v",
		instance.project.Labels.GroupId.Name, instance.project.GroupId,
		instance.project.Labels.ArtifactId.Name, instance.project.ArtifactId,
	)
}

func (instance *CleanupTask) getNamespaces() (result model.Namespaces, err error) {
	if v := instance.project.Claim.Namespaces; len(v) > 0 {
		return v, nil
	}
	resource := instance.client.Resource(schema.GroupVersionResource{
		Version:  "v1",
		Resource: "namespaces",
	})
	opts := metav1.ListOptions{}
	for {
		list, err := resource.List(opts)
		if err != nil {
			return nil, fmt.Errorf("cannot collect all existing namespaces: %w", err)
		}

		for _, candidate := range list.Items {
			result = append(result, model.Namespace(candidate.GetName()))
		}

		if v := list.GetContinue(); v != "" {
			opts.Continue = v
		} else {
			return nil, nil
		}
	}
}

type named map[model.Name]bool

func (instance *named) add(reference model.ObjectReference) {
	if instance == nil || *instance == nil {
		*instance = named{}
	}
	(*instance)[reference.Name] = true
}

func (instance named) has(reference model.ObjectReference) bool {
	if instance == nil {
		return false
	}
	return instance[reference.Name]
}

type namespaced map[model.Namespace]named

func (instance *namespaced) add(reference model.ObjectReference) {
	if instance == nil || *instance == nil {
		*instance = namespaced{}
	}
	v := (*instance)[reference.Namespace]
	v.add(reference)
	(*instance)[reference.Namespace] = v
}

func (instance namespaced) has(reference model.ObjectReference) bool {
	if instance == nil {
		return false
	}
	return instance[reference.Namespace].has(reference)
}

type gvked map[model.GroupVersionKind]namespaced

func (instance *gvked) add(reference model.ObjectReference) {
	if instance == nil || *instance == nil {
		*instance = gvked{}
	}
	v := (*instance)[reference.GroupVersionKind]
	v.add(reference)
	(*instance)[reference.GroupVersionKind] = v
}

func (instance gvked) has(reference model.ObjectReference) bool {
	if instance == nil {
		return false
	}
	return instance[reference.GroupVersionKind].has(reference)
}
