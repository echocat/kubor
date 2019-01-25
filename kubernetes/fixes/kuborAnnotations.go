package fixes

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func init() {
	registerUpdateFix(updateFixForAddingKuborAnnotations)
	registerCreateFix(createFixForAddingKuborAnnotations)
}

var ApplyKuborAnnotations = func(v1.Object, Project) error {
	return nil
}

func updateFixForAddingKuborAnnotations(project Project, _ unstructured.Unstructured, target *unstructured.Unstructured) error {
	return ApplyKuborAnnotations(target, project)
}

func createFixForAddingKuborAnnotations(project Project, target *unstructured.Unstructured) error {
	return ApplyKuborAnnotations(target, project)
}
