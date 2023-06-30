package transformation

import (
	"fmt"
	"github.com/echocat/kubor/model"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"reflect"
)

func init() {
	Default.MustRegisterUpdateFunc("service-preserve-cluster-ips", preserveServiceClusterIps)
	Default.MustRegisterUpdateFunc("service-preserve-health-check-node-port", preserveServiceHealthCheckNodePort)
	Default.MustRegisterUpdateFunc("service-preserve-node-ports", preserveServiceNodePorts)
}

var ServiceGvks = model.BuildGroupVersionKinds(v1.SchemeGroupVersion, &v1.Service{}).Build()

func preserveServiceClusterIps(_ *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured, _ *string) error {
	if !groupVersionKindMatches(&existing, target) {
		return nil
	}
	if !ServiceGvks.Contains(model.GroupVersionKind(target.GroupVersionKind())) {
		return nil
	}

	clusterIp, exist, err := unstructured.NestedString(existing.Object, "spec", "clusterIP")
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

func preserveServiceHealthCheckNodePort(_ *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured, _ *string) error {
	if !groupVersionKindMatches(&existing, target) {
		return nil
	}
	if !ServiceGvks.Contains(model.GroupVersionKind(target.GroupVersionKind())) {
		return nil
	}

	clusterIp, exist, err := unstructured.NestedInt64(existing.Object, "spec", "healthCheckNodePort")
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

	if _, exist, err := unstructured.NestedInt64(target.Object, "spec", "healthCheckNodePort"); err != nil {
		return err
	} else if exist {
		return nil
	}

	return unstructured.SetNestedField(target.Object, clusterIp, "spec", "healthCheckNodePort")
}

func preserveServiceNodePorts(_ *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured, _ *string) error {
	if !groupVersionKindMatches(&existing, target) {
		return nil
	}
	if !ServiceGvks.Contains(model.GroupVersionKind(target.GroupVersionKind())) {
		return nil
	}

	existingPorts, _, err := unstructured.NestedSlice(existing.Object, "spec", "ports")
	if err != nil {
		return err
	}
	targetPorts, _, err := unstructured.NestedSlice(target.Object, "spec", "ports")
	if err != nil {
		return err
	}
	if len(targetPorts) <= 0 {
		return nil
	}

	for i, pExistingPort := range existingPorts {
		name, existingPort, err := getNameOfUncheckedNamedMap(pExistingPort, i)
		if err != nil {
			return err
		}
		targetPort, targetIndex, err := findByNameOfUncheckedNamedMap(targetPorts, name)
		if err != nil {
			return err
		}
		if targetIndex >= 0 {
			existingNodePort, _, err := unstructured.NestedInt64(existingPort, "nodePort")
			if err != nil {
				return err
			}
			targetNodePort, _, err := unstructured.NestedInt64(targetPort, "nodePort")
			if err != nil {
				return err
			}

			if existingNodePort != 0 && targetNodePort == 0 {
				if err := unstructured.SetNestedField(targetPort, existingNodePort, "nodePort"); err != nil {
					return err
				}
				targetPorts[targetIndex] = targetPort
			}
		}
	}

	return unstructured.SetNestedSlice(target.Object, targetPorts, "spec", "ports")
}

func castUncheckedNamedMap(what interface{}, index int) (map[string]interface{}, error) {
	m, ok := what.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("existing port entry #%d should be of type map[string]interface{}, but got: %v", index, reflect.TypeOf(what))
	}
	return m, nil
}

func getNameOfUncheckedNamedMap(what interface{}, index int) (string, map[string]interface{}, error) {
	m, err := castUncheckedNamedMap(what, index)
	if err != nil {
		return "", nil, err
	}
	name, err := getNameOfNamedMap(m, index)
	return name, m, nil
}

func getNameOfNamedMap(what map[string]interface{}, index int) (string, error) {
	pv, ok := what["name"]
	if !ok {
		return "", fmt.Errorf("existing port entry #%d should contain field name, but it is not contained", index)
	}
	v, ok := pv.(string)
	if !ok {
		return "", fmt.Errorf("existing port entry #%d should contain field name of type string, but : %v", index, reflect.TypeOf(pv))
	}
	return v, nil
}

func matchesRequiredNameOfUncheckedNamedMap(what interface{}, index int, requiredName string) (bool, error) {
	m, ok := what.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("existing port entry #%d should be of type map[string]interface{}, but got: %v", index, reflect.TypeOf(what))
	}

	actualName, err := getNameOfNamedMap(m, index)
	if err != nil {
		// Ignore
		return false, nil
	}

	return requiredName == actualName, nil
}

func findByNameOfUncheckedNamedMap(candidates []interface{}, requiredName string) (map[string]interface{}, int, error) {
	for i, candidate := range candidates {
		matches, err := matchesRequiredNameOfUncheckedNamedMap(candidate, i, requiredName)
		if err != nil {
			return nil, -1, err
		}
		if matches {
			casted, err := castUncheckedNamedMap(candidate, i)
			if err != nil {
				return nil, -1, err
			}
			return casted, i, nil
		}
	}
	return nil, -1, nil
}
