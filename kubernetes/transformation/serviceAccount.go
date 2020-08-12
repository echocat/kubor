package transformation

import (
	"github.com/echocat/kubor/model"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func init() {
	Default.MustRegisterUpdateFunc("service-account-preserve-secrets", preserveServiceAccountSecrets)
}

var ServiceAccountGvks = model.BuildGroupVersionKinds(v1.SchemeGroupVersion, &v1.ServiceAccount{}).Build()

func preserveServiceAccountSecrets(_ *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured, _ string) error {
	if !groupVersionKindMatches(&existing, target) {
		return nil
	}
	if !ServiceAccountGvks.Contains(model.GroupVersionKind(target.GroupVersionKind())) {
		return nil
	}
	if v, isExplicitConfigured, _ := unstructured.NestedSlice(target.Object, "secrets"); isExplicitConfigured || len(v) > 0 {
		return nil
	}
	if v, isExplicitConfigured, _ := unstructured.NestedSlice(target.Object, "imagePullSecrets"); isExplicitConfigured || len(v) > 0 {
		return nil
	}

	if v, _, err := unstructured.NestedSlice(existing.Object, "secrets"); err != nil {
		return err
	} else if err := unstructured.SetNestedSlice(target.Object, v, "secrets"); err != nil {
		return err
	}
	if v, _, err := unstructured.NestedSlice(existing.Object, "imagePullSecrets"); err != nil {
		return err
	} else if err := unstructured.SetNestedSlice(target.Object, v, "imagePullSecrets"); err != nil {
		return err
	}

	return nil
}
