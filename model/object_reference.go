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
