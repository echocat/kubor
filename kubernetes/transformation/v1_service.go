package transformation

import (
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func init() {
	RegisterUpdateTransformationFunc(fixV1ServiceIfClusterIpIsAbsent)
}

func fixV1ServiceIfClusterIpIsAbsent(_ *model.Project, original unstructured.Unstructured, target *unstructured.Unstructured) error {
	if !groupVersionKindMatchesVersion(target, "v1") || !groupVersionKindMatchesKind(target, "service") {
		return nil
	}
	if !groupVersionKindMatches(&original, target) {
		return nil
	}

	clusterIp, exist, err := unstructured.NestedString(original.Object, "spec", "clusterIP")
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

	if _, exist, err := unstructured.NestedString(target.Object, "spec", "clusterIP"); err != nil {
		return err
	} else if exist {
		return nil
	}

	return unstructured.SetNestedField(target.Object, clusterIp, "spec", "clusterIP")
}
