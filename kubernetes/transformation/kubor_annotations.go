package transformation

import (
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func init() {
	RegisterUpdateTransformationFunc(func(project *model.Project, _ unstructured.Unstructured, target *unstructured.Unstructured) error {
		return ensureKuborAnnotations(project, target)
	})
	RegisterCreateTransformationFunc(ensureKuborLabels)
}

func ensureKuborAnnotations(project *model.Project, target *unstructured.Unstructured) error {
	pa := project.Annotations
	annotations := target.GetAnnotations()

	ensureKuborAnnotation(&annotations, pa.Stage)
	ensureKuborAnnotation(&annotations, pa.ApplyOn)
	ensureKuborAnnotation(&annotations, pa.WaitUntil)

	target.SetAnnotations(annotations)

	return nil
}

func ensureKuborAnnotation(annotations *map[string]string, annotation model.Annotation) {
	switch annotation.Action {
	case model.AnnotationActionDrop:
		delete(*annotations, annotation.Name.String())
	}
}
