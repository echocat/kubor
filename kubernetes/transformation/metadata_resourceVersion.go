package transformation

import (
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func init() {
	RegisterUpdateTransformationFunc(fixIfResourceVersionIsAbsent)
}

func fixIfResourceVersionIsAbsent(_ *model.Project, original unstructured.Unstructured, target *unstructured.Unstructured) error {
	if !groupVersionKindMatches(&original, target) {
		return nil
	}

	resourceVersion := original.GetResourceVersion()
	if resourceVersion == "" {
		return nil
	}

	target.SetResourceVersion(resourceVersion)

	return nil
}
