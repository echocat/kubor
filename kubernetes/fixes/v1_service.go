package fixes

import (
	"fmt"
	"github.com/echocat/kubor/common"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"reflect"
)

func init() {
	registerUpdateFix(updateFixForV1ServiceIfClusterIpIsAbsent)
}

func updateFixForV1ServiceIfClusterIpIsAbsent(original unstructured.Unstructured, target *unstructured.Unstructured) error {
	if !groupVersionKindMatchesVersion(target, "v1") || !groupVersionKindMatchesKind(target, "service") {
		return nil
	}
	if !groupVersionKindMatches(&original, target) {
		return nil
	}

	clusterIp := common.GetObjectPathValue(original.Object, "spec", "clusterIP")
	if clusterIp == nil {
		return nil
	}
	sClusterIp, ok := clusterIp.(string)
	if !ok || sClusterIp == "" {
		return nil
	}

	spec, ok := target.Object["spec"]
	if !ok {
		spec = map[string]interface{}{}
		target.Object["spec"] = spec
	}
	mSpec, ok := spec.(map[string]interface{})
	if !ok {
		return fmt.Errorf("'spec' property of target does already exists but is not of type map[string]interface{} it is %v", reflect.TypeOf(spec))
	}
	mSpec["clusterIP"] = sClusterIp

	return nil
}
