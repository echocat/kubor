package transformation

import (
	"github.com/echocat/kubor/model"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func init() {
	Default.MustRegisterUpdateFunc("service-preserve-cluster-ips", preserveServiceClusterIps)
	Default.MustRegisterUpdateFunc("service-preserve-health-check-node-port", preserveServiceHealthCheckNodePort)
	Default.MustRegisterUpdateFunc("service-preserve-node-ports", preserveServiceNodePorts)
}

var ServiceGvks = model.BuildGroupVersionKinds(v1.SchemeGroupVersion, &v1.Service{}).Build()

func preserveServiceClusterIps(_ *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured, _ string) error {
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

func preserveServiceHealthCheckNodePort(_ *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured, _ string) error {
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

func preserveServiceNodePorts(_ *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured, _ string) error {
	if !groupVersionKindMatches(&existing, target) {
		return nil
	}
	if !ServiceGvks.Contains(model.GroupVersionKind(target.GroupVersionKind())) {
		return nil
	}

	existingPorts, _, err := NestedNamedSliceAsMaps(existing.Object, "name", "spec", "ports")
	if err != nil {
		return err
	}
	targetPorts, _, err := NestedNamedSliceAsMaps(target.Object, "name", "spec", "ports")
	if err != nil {
		return err
	}
	if len(targetPorts) <= 0 {
		return nil
	}

	for name, existingPort := range existingPorts {
		if targetPort, ok := targetPorts[name]; ok {
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
				targetPorts[name] = targetPort
			}
		}
	}

	return SetNestedNamedMapsAsSlice(target.Object, "name", targetPorts, "spec", "ports")
}
