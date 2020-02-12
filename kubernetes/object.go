package kubernetes

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	rbacv1alpha1 "k8s.io/api/rbac/v1alpha1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"reflect"
)

var expectedNamespaceAbsentGvks = map[schema.GroupVersionKind]bool{
	corev1.SchemeGroupVersion.WithKind("namespace"): true,

	apiextensions.SchemeGroupVersion.WithKind("customresourcedefinition"):        true,
	apiextensionsv1.SchemeGroupVersion.WithKind("customresourcedefinition"):      true,
	apiextensionsv1beta1.SchemeGroupVersion.WithKind("customresourcedefinition"): true,

	rbacv1.SchemeGroupVersion.WithKind("clusterrole"):              true,
	rbacv1.SchemeGroupVersion.WithKind("clusterrolebinding"):       true,
	rbacv1beta1.SchemeGroupVersion.WithKind("clusterrole"):         true,
	rbacv1beta1.SchemeGroupVersion.WithKind("clusterrolebinding"):  true,
	rbacv1alpha1.SchemeGroupVersion.WithKind("clusterrole"):        true,
	rbacv1alpha1.SchemeGroupVersion.WithKind("clusterrolebinding"): true,
}

type ObjectValidator interface {
	IsNamespaced(what schema.GroupVersionKind) *bool
}

type ObjectInfo struct {
	Kind                 schema.GroupVersionKind
	Name                 string
	Namespace            string
	TypeMeta             metav1.TypeMeta
	GroupVersionResource schema.GroupVersionResource
}

func GetObjectInfo(object runtime.Object, by ObjectValidator) (ObjectInfo, error) {
	objk, ok := object.(schema.ObjectKind)
	if !ok {
		return ObjectInfo{}, fmt.Errorf("%v is not of type schema.ObjectKind", reflect.TypeOf(object))
	}
	objv, ok := object.(v1.Object)
	if !ok {
		return ObjectInfo{}, fmt.Errorf("%v is not of type v1.Object", reflect.TypeOf(object))
	}
	gvk := objk.GroupVersionKind()
	if gvk.Kind == "" {
		return ObjectInfo{}, fmt.Errorf("kind is not set or empty")
	}
	if gvk.Version == "" {
		return ObjectInfo{}, fmt.Errorf("apiVersion is not set or empty")
	}
	groupVersionResource, _ := meta.UnsafeGuessKindToResource(gvk)
	typeMeta := GroupVersionKindToTypeMeta(gvk)
	namespace := objv.GetNamespace()

	var namespaceExpectation bool
	if expectation := by.IsNamespaced(gvk); expectation != nil {
		namespaceExpectation = *expectation
	} else {
		namespaceExpectation = !expectedNamespaceAbsentGvks[NormalizeGroupVersionKind(gvk)]
	}
	if namespace == "" && namespaceExpectation {
		return ObjectInfo{}, fmt.Errorf("meta.namespace is not set or empty, but requird for %v", gvk)
	} else if namespace != "" && !namespaceExpectation {
		return ObjectInfo{}, fmt.Errorf("meta.namespace is set, but requird to be absent for %v", gvk)
	}
	name := objv.GetName()
	if name == "" {
		return ObjectInfo{}, fmt.Errorf("meta.name is not set or empty")
	}
	return ObjectInfo{
		Kind:                 gvk,
		Name:                 name,
		Namespace:            namespace,
		TypeMeta:             typeMeta,
		GroupVersionResource: groupVersionResource,
	}, nil
}

func (instance ObjectInfo) String() string {
	return fmt.Sprintf("%s/%s %s/%s", instance.TypeMeta.APIVersion, instance.TypeMeta.Kind, instance.Namespace, instance.Name)
}
