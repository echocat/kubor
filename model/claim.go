package model

import (
	"fmt"
	"github.com/echocat/kubor/template/functions"
)

type Claim struct {
	GroupVersionKinds GroupVersionKinds `yaml:"gvks,omitempty" json:"gvks,omitempty"`
	SourceNamespaces  []string          `yaml:"namespaces,omitempty" json:"namespaces,omitempty"`

	// Values set using implicitly.
	Namespaces Namespaces `yaml:"-" json:"-"`
}

var (
	DefaultClaimedGroupVersionKinds = DefaultGroupVersionKindRegistry.AsGroupVersionKinds()
)

func NewClaim() Claim {
	return Claim{
		GroupVersionKinds: DefaultClaimedGroupVersionKinds,
		SourceNamespaces:  []string{"{{.GroupId}}"},
	}
}

func (instance Claim) evaluate(context interface{}) (Claim, error) {
	fail := func(n string, err error) (Claim, error) {
		return Claim{}, fmt.Errorf("cannot handle namespace '%s': %w", n, err)
	}
	result := instance
	result.Namespaces = make(Namespaces, len(result.SourceNamespaces))
	for i, source := range result.SourceNamespaces {
		if tmpl, err := functions.DefaultTemplateFactory().New(source, source); err != nil {
			return fail(source, err)
		} else if rendered, err := tmpl.ExecuteToString(context); err != nil {
			return fail(source, err)
		} else if err := result.Namespaces[i].Set(rendered); err != nil {
			return fail(source, err)
		}
	}
	return result, nil
}

func (instance Claim) Validate(reference ObjectReference) error {
	if !instance.Namespaces.Contains(reference.Namespace) {
		return fmt.Errorf("is in namespace %v; but claimed: %v", reference.Namespace, instance.Namespaces)
	}
	if !instance.GroupVersionKinds.Contains(reference.GroupVersionKind) {
		return fmt.Errorf("is group version kind %v; but claimed: %v", reference.GroupVersionKind, instance.GroupVersionKinds)
	}
	return nil
}
