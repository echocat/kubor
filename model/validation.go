package model

import "k8s.io/apimachinery/pkg/runtime/schema"

type Validation struct {
	Schema SchemaValidation `json:"schema,omitempty" yaml:"schema,omitempty"`
}

type SchemaValidation struct {
	Ignored []GroupVersionKind `json:"ignored,omitempty" yaml:"ignored,omitempty"`
}

func (instance SchemaValidation) IsIgnored(what schema.GroupVersionKind) bool {
	for _, candidate := range instance.Ignored {
		if candidate.Group == what.Group &&
			candidate.Version == what.Version &&
			candidate.Kind == what.Kind {
			return true
		}
	}
	return false
}

type GroupVersionKind struct {
	Group   string `json:"group,omitempty" yaml:"group,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	Kind    string `json:"kind,omitempty" yaml:"kind,omitempty"`
}
