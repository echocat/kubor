package model

import (
	"strings"
)

type Scheme struct {
	Ignored    GroupVersionKinds            `json:"ignored,omitempty" yaml:"ignored,omitempty"`
	Namespaced []SchemaValidationNamespaced `json:"namespaced,omitempty" yaml:"namespaced,omitempty"`
}

func (instance Scheme) IsIgnored(what GroupVersionKind) bool {
	for candidate := range instance.Ignored {
		if strings.ToLower(candidate.Group) == strings.ToLower(what.Group) &&
			strings.ToLower(candidate.Version) == strings.ToLower(what.Version) &&
			strings.ToLower(candidate.Kind) == strings.ToLower(what.Kind) {
			return true
		}
	}
	return false
}

func (instance Scheme) IsNamespaced(what GroupVersionKind) *bool {
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

type SchemaValidationNamespaced struct {
	groupVersionKind `json:",inline" yaml:",inline"`

	Expectation bool `json:"expectation" yaml:"expectation"`
}
