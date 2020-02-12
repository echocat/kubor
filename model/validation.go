package model

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

type Validation struct {
	Schema SchemaValidation `json:"schema,omitempty" yaml:"schema,omitempty"`
}

type SchemaValidation struct {
	Ignored    []GroupVersionKind           `json:"ignored,omitempty" yaml:"ignored,omitempty"`
	Namespaced []SchemaValidationNamespaced `json:"namespaced,omitempty" yaml:"namespaced,omitempty"`
}

func (instance SchemaValidation) IsIgnored(what schema.GroupVersionKind) bool {
	for _, candidate := range instance.Ignored {
		if strings.ToLower(candidate.Group) == strings.ToLower(what.Group) &&
			strings.ToLower(candidate.Version) == strings.ToLower(what.Version) &&
			strings.ToLower(candidate.Kind) == strings.ToLower(what.Kind) {
			return true
		}
	}
	return false
}

func (instance SchemaValidation) IsNamespaced(what schema.GroupVersionKind) *bool {
	for _, candidate := range instance.Namespaced {
		if strings.ToLower(candidate.Group) == strings.ToLower(what.Group) &&
			strings.ToLower(candidate.Version) == strings.ToLower(what.Version) &&
			strings.ToLower(candidate.Kind) == strings.ToLower(what.Kind) {
			v := candidate.Expectation
			return &v
		}
	}
	return nil
}

type GroupVersionKind struct {
	Group   string `json:"group,omitempty" yaml:"group,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	Kind    string `json:"kind,omitempty" yaml:"kind,omitempty"`
}

type SchemaValidationNamespaced struct {
	GroupVersionKind `json:",inline" yaml:",inline"`

	Expectation bool `json:"expectation" yaml:"expectation"`
}
