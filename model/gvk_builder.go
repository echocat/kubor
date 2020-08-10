package model

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type GroupVersionKindsBuilder []groupVersionKindsBuilderEntry

func BuildGroupVersionKinds(gv schema.GroupVersion, obj runtime.Object) *GroupVersionKindsBuilder {
	result := &GroupVersionKindsBuilder{}
	return result.With(gv, obj)
}

func (instance *GroupVersionKindsBuilder) With(gv schema.GroupVersion, obj runtime.Object) *GroupVersionKindsBuilder {
	if instance == nil {
		*instance = GroupVersionKindsBuilder{}
	}
	*instance = append(*instance, groupVersionKindsBuilderEntry{gv, obj})
	return instance
}

func (instance GroupVersionKindsBuilder) Build() GroupVersionKinds {
	s := runtime.NewScheme()
	for _, candidate := range instance {
		s.AddKnownTypes(candidate.groupVersion, candidate.object)
	}
	return MapToGroupVersionKinds(s.AllKnownTypes())
}

type groupVersionKindsBuilderEntry struct {
	groupVersion schema.GroupVersion
	object       runtime.Object
}
