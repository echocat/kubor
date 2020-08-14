package transformation

import (
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strings"
)

const kuborAnnotationsName = model.TransformationName("apply-annotations")

func init() {
	t := ensureKuborAnnotations{}
	Default.MustRegisterUpdate(&t)
	Default.MustRegisterCreate(&t)
}

type ensureKuborAnnotations struct{}

func (instance *ensureKuborAnnotations) GetName() model.TransformationName {
	return kuborAnnotationsName
}

func (instance *ensureKuborAnnotations) GetPriority() int32 {
	return 1_000_000_000
}

func (instance *ensureKuborAnnotations) DefaultEnabled(target *unstructured.Unstructured) bool {
	return !NamespaceGvks.Contains(model.GroupVersionKind(target.GroupVersionKind()))
}

func (instance *ensureKuborAnnotations) TransformForUpdate(project *model.Project, _ unstructured.Unstructured, target *unstructured.Unstructured, argument *string) error {
	return instance.TransformForCreate(project, target, argument)
}

func (instance *ensureKuborAnnotations) TransformForCreate(project *model.Project, target *unstructured.Unstructured, _ *string) error {
	if err := instance.ensureOfPath(project, target, "metadata", "annotations"); err != nil {
		return err
	}
	if _, specTemplateExists, err := unstructured.NestedMap(target.Object, "spec", "template"); err != nil || !specTemplateExists {
		return err
	}
	if err := instance.ensureOfPath(project, target, "spec", "template", "metadata", "annotations"); err != nil {
		return err
	}
	return nil
}

func (instance *ensureKuborAnnotations) ensureOfPath(project *model.Project, target *unstructured.Unstructured, fields ...string) error {
	pa := project.Annotations
	annotations, _, err := unstructured.NestedStringMap(target.Object, fields...)
	if err != nil {
		return err
	}
	if annotations == nil {
		annotations = make(map[string]string)
	}

	instance.ensureAnnotation(&annotations, pa.Stage)
	instance.ensureAnnotation(&annotations, pa.ApplyOn)
	instance.ensureAnnotation(&annotations, pa.DryRunOn)
	instance.ensureAnnotation(&annotations, pa.WaitUntil)
	instance.ensureAnnotation(&annotations, pa.CleanupOn)
	instance.ensurePrefixedAnnotations(&annotations, pa.Transformations)

	return unstructured.SetNestedStringMap(target.Object, annotations, fields...)
}

func (instance *ensureKuborAnnotations) ensureAnnotation(annotations *map[string]string, annotation model.Annotation) {
	switch annotation.Action {
	case model.AnnotationActionDrop:
		delete(*annotations, string(annotation.Name))
	}
}

func (instance *ensureKuborAnnotations) ensurePrefixedAnnotations(annotations *map[string]string, annotation model.Annotation) {
	toHandle := map[string]bool{}
	for name := range *annotations {
		if strings.HasPrefix(name, string(annotation.Name)) {
			toHandle[name] = true
		}
	}
	for name := range toHandle {
		switch annotation.Action {
		case model.AnnotationActionDrop:
			delete(*annotations, name)
		}
	}
}
