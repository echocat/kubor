package transformation

import (
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strings"
)

const kuborAnnotationsTransformationName = model.TransformationName("apply-annotations")

func init() {
	Default.MustRegisterUpdateFunc(kuborAnnotationsTransformationName, ensureKuborAnnotationsOnUpdate)
	Default.MustRegisterCreateFunc(kuborAnnotationsTransformationName, ensureKuborAnnotations)
}

func ensureKuborAnnotationsOnUpdate(project *model.Project, _ unstructured.Unstructured, target *unstructured.Unstructured, argument *string) error {
	return ensureKuborAnnotations(project, target, argument)
}

func ensureKuborAnnotations(project *model.Project, target *unstructured.Unstructured, _ *string) error {
	if err := ensureKuborAnnotationsOfPath(project, target, "metadata", "annotations"); err != nil {
		return err
	}
	if _, specTemplateExists, err := unstructured.NestedMap(target.Object, "spec", "template"); err != nil || !specTemplateExists {
		return err
	}
	if err := ensureKuborAnnotationsOfPath(project, target, "spec", "template", "metadata", "annotations"); err != nil {
		return err
	}
	return nil
}

func ensureKuborAnnotationsOfPath(project *model.Project, target *unstructured.Unstructured, fields ...string) error {
	pa := project.Annotations
	annotations, _, err := unstructured.NestedStringMap(target.Object, fields...)
	if err != nil {
		return err
	}
	if annotations == nil {
		annotations = make(map[string]string)
	}

	ensureKuborAnnotation(&annotations, pa.Stage)
	ensureKuborAnnotation(&annotations, pa.ApplyOn)
	ensureKuborAnnotation(&annotations, pa.DryRunOn)
	ensureKuborAnnotation(&annotations, pa.WaitUntil)
	ensureKuborAnnotation(&annotations, pa.CleanupOn)
	ensureKuborPrefixedAnnotations(&annotations, pa.Transformations)

	return unstructured.SetNestedStringMap(target.Object, annotations, fields...)
}

func ensureKuborAnnotation(annotations *map[string]string, annotation model.Annotation) {
	switch annotation.Action {
	case model.AnnotationActionDrop:
		delete(*annotations, string(annotation.Name))
	}
}

func ensureKuborPrefixedAnnotations(annotations *map[string]string, annotation model.Annotation) {
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
