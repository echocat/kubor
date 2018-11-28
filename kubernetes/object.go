package kubernetes

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"reflect"
)

type ObjectInfo struct {
	Kind                 schema.GroupVersionKind
	Name                 string
	Namespace            string
	TypeMeta             metav1.TypeMeta
	GroupVersionResource schema.GroupVersionResource
}

func GetObjectInfo(object v1.Object) (ObjectInfo, error) {
	objk, ok := object.(schema.ObjectKind)
	if !ok {
		return ObjectInfo{}, fmt.Errorf("%v is not of type schema.ObjectKind", reflect.TypeOf(object))
	}
	kind := objk.GroupVersionKind()
	if kind.Kind == "" {
		return ObjectInfo{}, fmt.Errorf("kind is not set or empty")
	}
	if kind.Version == "" {
		return ObjectInfo{}, fmt.Errorf("apiVersion is not set or empty")
	}
	groupVersionResource, _ := meta.UnsafeGuessKindToResource(kind)
	typeMeta := GroupVersionKindToTypeMeta(kind)
	namespace := object.GetNamespace()
	if namespace == "" {
		return ObjectInfo{}, fmt.Errorf("meta.namespace is not set or empty")
	}
	name := object.GetName()
	if name == "" {
		return ObjectInfo{}, fmt.Errorf("meta.name is not set or empty")
	}
	return ObjectInfo{
		Kind:                 kind,
		Name:                 name,
		Namespace:            namespace,
		TypeMeta:             typeMeta,
		GroupVersionResource: groupVersionResource,
	}, nil
}

func (instance ObjectInfo) String() string {
	return fmt.Sprintf("%s/%s %s/%s", instance.TypeMeta.APIVersion, instance.TypeMeta.Kind, instance.Namespace, instance.Name)
}
