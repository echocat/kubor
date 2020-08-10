package model

import (
	"encoding/json"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"reflect"
	"strings"
)

type GroupVersionKind schema.GroupVersionKind

func (instance GroupVersionKind) Normalize() GroupVersionKind {
	return GroupVersionKind{
		Group:   strings.ToLower(instance.Group),
		Version: strings.ToLower(instance.Version),
		Kind:    strings.ToLower(instance.Kind),
	}
}

func (instance GroupVersionKind) String() string {
	result := instance.Kind
	if v := instance.Version; v != "" {
		result = v + "/" + result
	}
	if v := instance.Group; v != "" {
		result = v + "/" + result
	}
	return result
}

func (instance GroupVersionKind) Bare() schema.GroupVersionKind {
	return schema.GroupVersionKind(instance)
}

func (instance GroupVersionKind) GroupKind() schema.GroupKind {
	return schema.GroupVersionKind(instance).GroupKind()
}

func (instance GroupVersionKind) GroupVersion() schema.GroupVersion {
	return schema.GroupVersionKind(instance).GroupVersion()
}

func (instance GroupVersionKind) GuessToResource() (plural, singular schema.GroupVersionResource) {
	return meta.UnsafeGuessKindToResource(schema.GroupVersionKind(instance))
}

type groupVersionKind struct {
	Group   string `json:"group,omitempty" yaml:"group,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	Kind    string `json:"kind,omitempty" yaml:"kind,omitempty"`
}

func (instance groupVersionKind) get() GroupVersionKind {
	return GroupVersionKind{
		Group:   instance.Group,
		Version: instance.Version,
		Kind:    instance.Kind,
	}.Normalize()
}

func (instance *groupVersionKind) set(v GroupVersionKind) {
	instance.Group = v.Group
	instance.Version = v.Version
	instance.Kind = v.Kind
}

type groupVersionKinds []groupVersionKind

func (instance groupVersionKinds) get() GroupVersionKinds {
	result := make(GroupVersionKinds, len(instance))
	for _, v := range instance {
		result[v.get()] = true
	}
	return result
}

type GroupVersionKinds map[GroupVersionKind]bool

func (instance GroupVersionKinds) Contains(v GroupVersionKind) bool {
	if len(instance) == 0 {
		return true
	}
	v = v.Normalize()
	for candidate := range instance {
		if v == candidate {
			return true
		}
	}
	return false
}

func (instance GroupVersionKinds) Strings() []string {
	result := make([]string, len(instance))
	var i int
	for v := range instance {
		result[i] = v.String()
		i++
	}
	return result
}

func (instance GroupVersionKinds) String() string {
	return strings.Join(instance.Strings(), ",")
}

func (instance GroupVersionKinds) get() groupVersionKinds {
	result := make(groupVersionKinds, len(instance))
	var i int
	for v := range instance {
		var nv groupVersionKind
		nv.set(v)
		result[i] = nv
		i++
	}
	return result
}

func (instance *GroupVersionKinds) UnmarshalJSON(b []byte) error {
	var buf groupVersionKinds
	if err := json.Unmarshal(b, &buf); err != nil {
		return err
	}
	*instance = buf.get()
	return nil
}

func (instance GroupVersionKinds) MarshalJSON() ([]byte, error) {
	return json.Marshal(instance.get())
}

func (instance *GroupVersionKinds) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var buf groupVersionKinds
	if err := unmarshal(&buf); err != nil {
		return err
	}
	*instance = buf.get()
	return nil
}

func (instance GroupVersionKinds) MarshalYAML() (interface{}, error) {
	return instance.get(), nil
}

func MapToGroupVersionKinds(in map[schema.GroupVersionKind]reflect.Type) GroupVersionKinds {
	result := make(GroupVersionKinds, len(in))
	for v := range in {
		result[GroupVersionKind(v).Normalize()] = true
	}
	return result
}
