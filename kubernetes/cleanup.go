package kubernetes

import (
	"context"
	"fmt"
	"github.com/echocat/kubor/model"
	"github.com/echocat/slf4g"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"time"
)

type CleanupTask struct {
	project *model.Project
	keep    gvked
	client  dynamic.Interface
	mode    CleanupMode
}

func NewCleanupTask(project *model.Project, client dynamic.Interface, mode CleanupMode) (CleanupTask, error) {
	return CleanupTask{
		project: project,
		client:  client,
		mode:    mode,
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
	l := log.With("namespace", namespace).
		With("mode", instance.mode)

	start := time.Now()

	defer func() {
		l = l.With("duration", time.Now().Sub(start))
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

	handledGvks := model.GroupVersionKinds{}
	for gvk := range instance.project.Claim.GroupVersionKinds {
		respect := true
		for twin := range model.DefaultGroupVersionKindRegistry.GetTwins(gvk) {
			if handledGvks[twin] {
				respect = false
			}
		}

		if respect {
			if foundAtLeastOne, err := instance.executeFor(l, namespace, gvk); err != nil {
				return err
			} else if foundAtLeastOne {
				handledGvks[gvk] = true
			}
		}
	}

	return nil
}

func (instance *CleanupTask) executeFor(l log.Logger, namespace model.Namespace, gvk model.GroupVersionKind) (foundAtLeastOne bool, err error) {
	l = l.With("gvk", gvk)

	start := time.Now()

	defer func() {
		l = l.With("duration", time.Now().Sub(start))
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
		list, err := resource.List(context.Background(), opts)
		if err != nil {
			if as, ok := err.(errors.APIStatus); ok && as.Status().Code == 404 {
				return false, nil
			}
			return false, fmt.Errorf("cannot collect existing elements of type %v: %w", gvr, err)
		}

		for _, candidate := range list.Items {
			foundAtLeastOne = true
			reference, err := GetObjectReference(&candidate, instance.project.Scheme)
			if err != nil {
				l.WithError(err).
					With("reference", fmt.Sprintf("%v %v/%v", candidate.GroupVersionKind(), candidate.GetName(), candidate.GetNamespace())).
					Warn("Cannot evaluate %v. Skipping it...", instance.mode.AffectedDescription(false, false))
				continue
			}

			if instance.shouldBeKept(reference) {
				l.With("reference", reference).
					Trace("%v %v is part of the deployment and will be kept.", instance.mode.AffectedDescription(false, true), reference)
				continue
			}

			if instance.hasOwner(&candidate) {
				l.With("reference", reference).
					Trace("%v %v has an owner and will therefore be kept.", instance.mode.AffectedDescription(false, true), reference)
				continue
			}

			if allowedToBeDeleted, err := instance.isAllowedToBeDeleted(&candidate); err != nil {
				return false, err
			} else if !allowedToBeDeleted {
				l.With("reference", reference).
					Trace("%v %v is not allowed to be deleted in mode %v and will therefore be kept.",
						instance.mode.AffectedDescription(false, true), reference, instance.mode)
				continue
			}

			if err := instance.delete(resource, reference); err != nil {
				return false, err
			}
		}

		if v := list.GetContinue(); v != "" {
			opts.Continue = v
		} else {
			return foundAtLeastOne, nil
		}
	}
}

func (instance *CleanupTask) shouldBeKept(reference model.ObjectReference) bool {
	if len(instance.keep) == 0 {
		return false
	}
	if instance.keep.has(reference) {
		return true
	}
	for _, twin := range reference.AllTwinsBy(model.DefaultGroupVersionKindRegistry) {
		if instance.keep.has(twin) {
			return true
		}
	}
	return false
}

func (instance *CleanupTask) delete(resource dynamic.ResourceInterface, reference model.ObjectReference) (err error) {
	start := time.Now()
	l := log.
		With("action", "delete")

	defer func() {
		ld := l.With("duration", time.Now().Sub(start))
		if err != nil {
			ldd := ld.
				WithError(err).
				With("status", "failed")
			if ldd.IsDebugEnabled() {
				ldd.Error("Deleting %v %v... FAILED!", instance.mode.AffectedDescription(false, false), reference)
			} else {
				ldd.Error("Was not able to delete %v %v.", instance.mode.AffectedDescription(false, false), reference)
			}
		} else {
			ldd := ld.With("status", "success")
			if ldd.IsDebugEnabled() {
				ldd.Info("Deleting %v %v... DONE!", instance.mode.AffectedDescription(false, false), reference)
			} else {
				ldd.Info("%v %v deleted.", instance.mode.AffectedDescription(false, true), reference)
			}
		}
	}()

	l.Debug("Deleting %v %v...", instance.mode.AffectedDescription(false, false), reference)

	dp := metav1.DeletePropagationForeground
	if err := resource.Delete(context.Background(), reference.Name.String(), metav1.DeleteOptions{
		PropagationPolicy: &dp,
	}); err != nil {
		return fmt.Errorf("cannot delete %v: %w", instance.mode.AffectedDescription(false, false), err)
	}

	return nil
}

func (instance *CleanupTask) hasOwner(target *unstructured.Unstructured) bool {
	result := len(target.GetOwnerReferences()) > 0
	return result
}

func (instance *CleanupTask) isAllowedToBeDeleted(target *unstructured.Unstructured) (bool, error) {
	rule, err := instance.project.Annotations.GetCleanupOn(target)
	if err != nil {
		return false, err
	}
	switch instance.mode {
	case CleanupModeOrphans:
		return rule.OnOrphaned(), nil
	case CleanupModeDelete:
		return rule.OnDelete(), nil
	default:
		return false, nil
	}
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
		list, err := resource.List(context.Background(), opts)
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

type CleanupMode uint8

const (
	CleanupModeOrphans = CleanupMode(0)
	CleanupModeDelete  = CleanupMode(1)
)

func (instance CleanupMode) String() string {
	switch instance {
	case CleanupModeOrphans:
		return "orphans"
	case CleanupModeDelete:
		return "everything"
	default:
		return fmt.Sprintf("unknown-%d", instance)
	}
}

func (instance CleanupMode) AffectedDescription(plural bool, capitalize bool) (result string) {
	defer func() {
		if capitalize {
			result = cases.Title(language.AmericanEnglish).String(result)
		}
	}()

	switch instance {
	case CleanupModeOrphans:
		if plural {
			return "orphans"
		}
		return "orphan"
	default:
		if plural {
			return "elements"
		}
		return "element"
	}
}
