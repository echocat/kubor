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
	Updates map[model.TransformationName]Update
	Creates map[model.TransformationName]Create
}

type Update interface {
	TransformForUpdate(p *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured, argument *string) error
	DefaultEnabled(target *unstructured.Unstructured) bool
}

type Create interface {
	TransformForCreate(p *model.Project, target *unstructured.Unstructured, argument *string) error
	DefaultEnabled(target *unstructured.Unstructured) bool
}

type UpdateFunc func(p *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured, argument *string) error

func (instance UpdateFunc) TransformForUpdate(p *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured, argument *string) error {
	return instance(p, existing, target, argument)
}

func (instance UpdateFunc) DefaultEnabled(*unstructured.Unstructured) bool {
	return true
}

type CreateFunc func(p *model.Project, target *unstructured.Unstructured, argument *string) error

func (instance CreateFunc) TransformForCreate(p *model.Project, target *unstructured.Unstructured, argument *string) error {
	return instance(p, target, argument)
}

func (instance CreateFunc) DefaultEnabled(*unstructured.Unstructured) bool {
	return true
}

func (instance *Transformations) RegisterUpdate(name model.TransformationName, v Update) error {
	if _, err := name.MarshalText(); err != nil {
		return err
	}
	if instance == nil {
		*instance = Transformations{}
	}
	if instance.Updates == nil {
		instance.Updates = map[model.TransformationName]Update{}
	}
	instance.Updates[name] = v
	return nil
}

func (instance *Transformations) MustRegisterUpdate(name model.TransformationName, v Update) {
	if err := instance.RegisterUpdate(name, v); err != nil {
		panic(err)
	}
}

func (instance *Transformations) RegisterUpdateFunc(name model.TransformationName, v UpdateFunc) error {
	return instance.RegisterUpdate(name, v)
}

func (instance *Transformations) MustRegisterUpdateFunc(name model.TransformationName, v UpdateFunc) {
	if err := instance.RegisterUpdateFunc(name, v); err != nil {
		panic(err)
	}
}

func (instance *Transformations) RegisterCreate(name model.TransformationName, v Create) error {
	if _, err := name.MarshalText(); err != nil {
		return err
	}
	if instance == nil {
		*instance = Transformations{}
	}
	if instance.Creates == nil {
		instance.Creates = map[model.TransformationName]Create{}
	}
	instance.Creates[name] = v
	return nil
}

func (instance *Transformations) MustRegisterCreate(name model.TransformationName, v Create) {
	if err := instance.RegisterCreate(name, v); err != nil {
		panic(err)
	}
}

func (instance *Transformations) RegisterCreateFunc(name model.TransformationName, v CreateFunc) error {
	return instance.RegisterCreate(name, v)
}

func (instance *Transformations) MustRegisterCreateFunc(name model.TransformationName, v CreateFunc) {
	if err := instance.RegisterCreateFunc(name, v); err != nil {
		panic(err)
	}
}

func (instance Transformations) TransformForUpdate(p *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured) error {
	if v := instance.Updates; v != nil {
		for name, transformation := range v {
			v, err := p.GetTransformation(target, name)
			if err != nil {
				return fmt.Errorf("cannot evaluate annotations for transformation %v: %w", name, err)
			}
			if !v.IsEnabled(transformation.DefaultEnabled(target)) {
				continue
			}
			if err := transformation.TransformForUpdate(p, existing, target, v.Argument); err != nil {
				return fmt.Errorf("cannot evaluate transformation %v with argument %s: %w", name, v.ArgumentAsString(), err)
			}
		}
	}
	return nil
}

func (instance Transformations) TransformForCreate(p *model.Project, target *unstructured.Unstructured) error {
	if v := instance.Creates; v != nil {
		for name, transformation := range v {
			v, err := p.GetTransformation(target, name)
			if err != nil {
				return fmt.Errorf("cannot evaluate annotations for transformation %v: %w", name, err)
			}
			if !v.IsEnabled(transformation.DefaultEnabled(target)) {
				continue
			}
			if err := transformation.TransformForCreate(p, target, v.Argument); err != nil {
				return fmt.Errorf("cannot evaluate transformation %v with argument %v: %w", name, v.ArgumentAsString(), err)
			}
		}
	}
	return nil
}

func (instance Transformations) Clone() (result Transformations) {
	result.Updates = make(map[model.TransformationName]Update, len(instance.Updates))
	result.Creates = make(map[model.TransformationName]Create, len(instance.Creates))

	for name, candidate := range instance.Updates {
		result.Updates[name] = candidate
	}
	for name, candidate := range instance.Creates {
		result.Creates[name] = candidate
	}

	return
}
