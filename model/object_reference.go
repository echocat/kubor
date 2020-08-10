package model

import (
	"fmt"
)

type ObjectReference struct {
	GroupVersionKind
	Name      Name
	Namespace Namespace
}

func (instance ObjectReference) String() string {
	return fmt.Sprintf("%v %s/%s", instance.GroupVersionKind, instance.Namespace, instance.Name)
}

func (instance ObjectReference) AllTwinsBy(registry *GroupVersionKindRegistry) (result []ObjectReference) {
	gvkTwins := registry.GetTwins(instance.GroupVersionKind)
	result = make([]ObjectReference, len(gvkTwins))
	var i int
	for gvkTwin := range gvkTwins {
		result[i] = ObjectReference{
			GroupVersionKind: gvkTwin,
			Name:             instance.Name,
			Namespace:        instance.Namespace,
		}
		i++
	}
	return
}
