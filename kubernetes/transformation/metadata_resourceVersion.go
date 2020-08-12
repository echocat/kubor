package transformation

import (
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func init() {
	Default.RegisterUpdateFunc("preserve-ResourceVersion", preserveResourceVersion)
}

func preserveResourceVersion(_ *model.Project, original unstructured.Unstructured, target *unstructured.Unstructured, _ string) error {
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
