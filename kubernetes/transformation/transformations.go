package transformation

import (
	"fmt"
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	Default = Transformations{}
)

type Transformations struct {
	Updates map[Name]Update
	Creates map[Name]Create
}

type Update interface {
	TransformForUpdate(p *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured, argument string) error
}

type Create interface {
	TransformForCreate(p *model.Project, target *unstructured.Unstructured, argument string) error
}

type UpdateFunc func(p *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured, argument string) error

func (instance UpdateFunc) TransformForUpdate(p *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured, argument string) error {
	return instance(p, existing, target, argument)
}

type CreateFunc func(p *model.Project, target *unstructured.Unstructured, argument string) error

func (instance CreateFunc) TransformForCreate(p *model.Project, target *unstructured.Unstructured, argument string) error {
	return instance(p, target, argument)
}

func (instance *Transformations) RegisterUpdate(name Name, v Update) Update {
	if instance == nil {
		*instance = Transformations{}
	}
	if instance.Updates == nil {
		instance.Updates = map[Name]Update{}
	}
	instance.Updates[name] = v
	return v
}

func (instance *Transformations) RegisterUpdateFunc(name Name, v UpdateFunc) {
	instance.RegisterUpdate(name, v)
}

func (instance *Transformations) RegisterCreate(name Name, v Create) Create {
	if instance == nil {
		*instance = Transformations{}
	}
	if instance.Creates == nil {
		instance.Creates = map[Name]Create{}
	}
	instance.Creates[name] = v
	return v
}

func (instance *Transformations) RegisterCreateFunc(name Name, v CreateFunc) {
	instance.RegisterCreate(name, v)
}

func (instance Transformations) TransformForUpdate(p *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured) error {
	if v := instance.Updates; v != nil {
		for name, transformation := range v {
			enabled, argument, err := p.Annotations.GetTransformationState(target, name.String(), true)
			if err != nil {
				return fmt.Errorf("cannot evaluate annotations for transformation %v: %w", name, err)
			}
			if !enabled {
				continue
			}
			if err := transformation.TransformForUpdate(p, existing, target, argument); err != nil {
				return fmt.Errorf("cannot evaluate transformation %v with argument %q: %w", name, argument, err)
			}
		}
	}
	return nil
}

func (instance Transformations) TransformForCreate(p *model.Project, target *unstructured.Unstructured) error {
	if v := instance.Creates; v != nil {
		for name, transformation := range v {
			enabled, argument, err := p.Annotations.GetTransformationState(target, name.String(), true)
			if err != nil {
				return fmt.Errorf("cannot evaluate annotations for transformation %v: %w", name, err)
			}
			if !enabled {
				continue
			}
			if err := transformation.TransformForCreate(p, target, argument); err != nil {
				return fmt.Errorf("cannot evaluate transformation %v with argument %q: %w", name, argument, err)
			}
		}
	}
	return nil
}

func (instance Transformations) Clone() (result Transformations) {
	result.Updates = make(map[Name]Update, len(instance.Updates))
	result.Creates = make(map[Name]Create, len(instance.Creates))

	for name, candidate := range instance.Updates {
		result.Updates[name] = candidate
	}
	for name, candidate := range instance.Creates {
		result.Creates[name] = candidate
	}

	return
}
