package model

import (
	"fmt"
	"github.com/levertonai/kubor/common"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"reflect"
	"regexp"
	"strings"
)

type VersionKind struct {
	GroupVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind         string `yaml:"kind" json:"kind"`
}

func (instance VersionKind) Matches(gvk schema.GroupVersionKind) bool {
	return strings.ToLower(gvk.Group) == strings.ToLower(instance.Group()) &&
		strings.ToLower(gvk.Version) == strings.ToLower(instance.Version()) &&
		strings.ToLower(gvk.Kind) == strings.ToLower(instance.Kind)
}

func (instance VersionKind) Group() string {
	parts := strings.SplitN(instance.GroupVersion, "/", 2)
	if len(parts) == 1 {
		return ""
	}
	return parts[0]
}

func (instance VersionKind) Version() string {
	parts := strings.SplitN(instance.GroupVersion, "/", 2)
	if len(parts) == 1 {
		return instance.GroupVersion
	}
	return parts[1]
}

func (instance VersionKind) eval(project Project) (result VersionKind, err error) {
	if result.GroupVersion, err = evaluateTemplate("claim.kind.apiVersion", instance.GroupVersion, project); err != nil {
		return
	}
	if result.Kind, err = evaluateTemplate("claim.kind.kind", instance.Kind, project); err != nil {
		return
	}
	return
}

type VersionKinds []VersionKind

func (instance VersionKinds) eval(project Project) (result VersionKinds, err error) {
	result = make(VersionKinds, len(instance))
	for i, part := range instance {
		if result[i], err = part.eval(project); err != nil {
			return
		}
	}
	return
}

func (instance VersionKinds) Matches(gvk schema.GroupVersionKind) bool {
	for _, candidate := range instance {
		if candidate.Matches(gvk) {
			return true
		}
	}
	return false
}

type Namespace string

func (instance Namespace) eval(project Project) (Namespace, error) {
	if strResult, err := evaluateTemplate("claim.namespace", string(instance), project); err != nil {
		return Namespace(""), err
	} else {
		return Namespace(strResult), nil
	}
}

func (instance Namespace) Matches(namespace string) bool {
	return strings.ToLower(string(instance)) == strings.ToLower(namespace)
}

type Namespaces []Namespace

func (instance Namespaces) eval(project Project) (result Namespaces, err error) {
	result = make(Namespaces, len(instance))
	for i, part := range instance {
		if result[i], err = part.eval(project); err != nil {
			return
		}
	}
	return
}

func (instance Namespaces) Matches(namespace string) bool {
	for _, candidate := range instance {
		if candidate.Matches(namespace) {
			return true
		}
	}
	return false
}

type Name struct {
	source  string
	pattern *regexp.Regexp
}

func (instance *Name) UnmarshalYAML(unmarshal func(interface{}) error) error {
	instance.pattern = nil
	return unmarshal(&instance.source)
}

func (instance Name) MarshalYAML() (interface{}, error) {
	return instance.source, nil
}

func (instance Name) eval(project Project) (Name, error) {
	if source, err := evaluateTemplate("claim.name", instance.source, project); err != nil {
		return Name{}, err
	} else if pattern, err := regexp.Compile("^" + source + "$"); err != nil {
		return Name{}, nil
	} else {
		return Name{
			source:  source,
			pattern: pattern,
		}, nil
	}
}

func (instance Name) Matches(name string) (bool, error) {
	pattern := instance.pattern
	if pattern == nil {
		return false, nil
	} else {
		return pattern.MatchString(name), nil
	}
}

type Names []Name

func (instance Names) eval(project Project) (result Names, err error) {
	result = make(Names, len(instance))
	for i, part := range instance {
		if result[i], err = part.eval(project); err != nil {
			return
		}
	}
	return
}

func (instance Names) Matches(name string) (bool, error) {
	for _, candidate := range instance {
		if matches, err := candidate.Matches(name); err != nil {
			return false, err
		} else if matches {
			return true, nil
		}
	}
	return false, nil
}

type Claim struct {
	Kinds      VersionKinds `yaml:"kinds,omitempty" json:"kinds,omitempty"`
	Namespaces Namespaces   `yaml:"namespaces,omitempty" json:"namespaces,omitempty"`
	Names      Names        `yaml:"names,omitempty" json:"names,omitempty"`
}

func (instance Claim) Matches(object runtime.Object) (bool, error) {
	if !instance.Kinds.Matches(object.GetObjectKind().GroupVersionKind()) {
		return false, nil
	}
	objv, ok := object.(v1.Object)
	if !ok {
		return false, fmt.Errorf("%v is not of type v1.Object", reflect.TypeOf(object))
	}
	if matches := instance.Namespaces.Matches(objv.GetNamespace()); !matches {
		return false, nil
	}
	if matches, err := instance.Names.Matches(objv.GetName()); err != nil || !matches {
		return false, err
	}
	return true, nil
}

func (instance Claim) eval(project Project) (result Claim, err error) {
	if result.Kinds, err = instance.Kinds.eval(project); err != nil {
		return
	}
	if result.Namespaces, err = instance.Namespaces.eval(project); err != nil {
		return
	}
	if result.Names, err = instance.Names.eval(project); err != nil {
		return
	}
	return
}

type Claims map[string]Claim

func (instance Claims) Matches(object runtime.Object) (bool, error) {
	for _, candidate := range instance {
		if matches, err := candidate.Matches(object); err != nil {
			return false, err
		} else if matches {
			return true, nil
		}
	}
	return false, nil
}

type ClaimDefinition struct {
	On    common.EvaluatingPredicate `yaml:"on,omitempty" json:"on,omitempty"`
	Claim Claim                      `yaml:",inline" json:",inline"`
}

type ClaimDefinitions map[string]ClaimDefinition

func (instance ClaimDefinitions) ToClaims(project Project) (result Claims, _ error) {
	result = Claims{}
	for name, definition := range instance {
		if matches, err := definition.On.Matches(project); err != nil {
			return Claims{}, err
		} else if matches {
			if claim, err := definition.Claim.eval(project); err != nil {
				return Claims{}, err
			} else {
				result[name] = claim
			}
		}
	}
	return result, nil
}
