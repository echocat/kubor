package kubernetes

import (
	"fmt"
	"github.com/echocat/kubor/model"
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

var expectedNamespaceAbsentGvks = func() model.GroupVersionKinds {
	s := runtime.NewScheme()
	s.AddKnownTypes(corev1.SchemeGroupVersion,
		&corev1.Namespace{},
	)
	s.AddKnownTypes(apiextensions.SchemeGroupVersion,
		&apiextensions.CustomResourceDefinition{},
	)
	s.AddKnownTypes(apiextensionsv1.SchemeGroupVersion,
		&apiextensionsv1.CustomResourceDefinition{},
	)
	s.AddKnownTypes(apiextensionsv1beta1.SchemeGroupVersion,
		&apiextensionsv1beta1.CustomResourceDefinition{},
	)
	s.AddKnownTypes(rbacv1.SchemeGroupVersion,
		&rbacv1.ClusterRole{},
		&rbacv1.ClusterRoleBinding{},
	)
	s.AddKnownTypes(rbacv1beta1.SchemeGroupVersion,
		&rbacv1beta1.ClusterRole{},
		&rbacv1beta1.ClusterRoleBinding{},
	)
	s.AddKnownTypes(rbacv1alpha1.SchemeGroupVersion,
		&rbacv1alpha1.ClusterRole{},
		&rbacv1alpha1.ClusterRoleBinding{},
	)
	return model.MapToGroupVersionKinds(s.AllKnownTypes())
}()

type ObjectValidator interface {
	IsNamespaced(what model.GroupVersionKind) *bool
}

type ObjectInfo struct {
	model.ObjectReference
	TypeMeta metav1.TypeMeta
	Resource schema.GroupVersionResource
}

func GetObjectInfo(object runtime.Object, by ObjectValidator) (ObjectInfo, error) {
	reference, err := GetObjectReference(object, by)
	if err != nil {
		return ObjectInfo{}, err
	}
	groupVersionResource, _ := meta.UnsafeGuessKindToResource(reference.GroupVersionKind.Bare())
	typeMeta := GroupVersionKindToTypeMeta(reference.GroupVersionKind)

	return ObjectInfo{
		ObjectReference: reference,
		TypeMeta:        typeMeta,
		Resource:        groupVersionResource,
	}, nil
}

func (instance ObjectInfo) String() string {
	return instance.ObjectReference.String()
}

func GetObjectReference(object runtime.Object, by ObjectValidator) (model.ObjectReference, error) {
	objk, ok := object.(schema.ObjectKind)
	if !ok {
		return model.ObjectReference{}, fmt.Errorf("%v is not of type schema.ObjectKind", reflect.TypeOf(object))
	}
	objv, ok := object.(v1.Object)
	if !ok {
		return model.ObjectReference{}, fmt.Errorf("%v is not of type v1.Object", reflect.TypeOf(object))
	}
	gvk := model.GroupVersionKind(objk.GroupVersionKind()).Normalize()
	if gvk.Kind == "" {
		return model.ObjectReference{}, fmt.Errorf("kind is not set or empty")
	}
	if gvk.Version == "" {
		return model.ObjectReference{}, fmt.Errorf("apiVersion is not set or empty")
	}
	namespace := model.Namespace(objv.GetNamespace())

	var namespaceExpectation bool
	if expectation := by.IsNamespaced(gvk); expectation != nil {
		namespaceExpectation = *expectation
	} else {
		_, found := expectedNamespaceAbsentGvks[gvk]
		namespaceExpectation = !found
	}
	if namespace == "" && namespaceExpectation {
		return model.ObjectReference{}, fmt.Errorf("meta.namespace is not set or empty, but requird for %v", gvk)
	} else if namespace != "" && !namespaceExpectation {
		return model.ObjectReference{}, fmt.Errorf("meta.namespace is set, but requird to be absent for %v", gvk)
	} else if _, err := namespace.MarshalText(); err != nil {
		return model.ObjectReference{}, fmt.Errorf("illegal meta.namespace for %v: %w", gvk, err)
	}
	pname := objv.GetName()
	if pname == "" {
		return model.ObjectReference{}, fmt.Errorf("meta.name is not set or empty")
	}
	var name model.Name
	if err := name.Set(pname); err != nil {
		return model.ObjectReference{}, fmt.Errorf("illegal meta.name: %w", err)
	}

	return model.ObjectReference{
		GroupVersionKind: gvk,
		Name:             name,
		Namespace:        namespace,
	}, nil
}
