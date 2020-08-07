package transformation

import (
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func init() {
	RegisterUpdateTransformationFunc(func(project *model.Project, _ unstructured.Unstructured, target *unstructured.Unstructured) error {
		return ensureKuborAnnotations(project, target)
	})
	RegisterCreateTransformationFunc(ensureKuborAnnotations)
}

func ensureKuborAnnotations(project *model.Project, target *unstructured.Unstructured) error {
	if err := ensureKuborAnnotationsOfPath(project, target, "metadata", "annotations"); err != nil {
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

	return unstructured.SetNestedStringMap(target.Object, annotations, fields...)
}

func ensureKuborAnnotation(annotations *map[string]string, annotation model.Annotation) {
	switch annotation.Action {
	case model.AnnotationActionDrop:
		delete(*annotations, annotation.Name.String())
	}
}
