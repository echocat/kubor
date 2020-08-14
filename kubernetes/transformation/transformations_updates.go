package transformation

import (
	"fmt"
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sort"
)

type Update interface {
	Transformation
	TransformForUpdate(p *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured, argument *string) error
}

type Updates []Update

func (instance Updates) TransformForUpdate(p *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured) error {
	for _, transformation := range instance {
		name := transformation.GetName()
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
	return nil
}

func (instance Updates) Len() int      { return len(instance) }
func (instance Updates) Swap(i, j int) { instance[i], instance[j] = instance[j], instance[i] }
func (instance Updates) Less(i, j int) bool {
	vI := instance[i]
	vJ := instance[j]
	if vJ == nil {
		return false
	}
	if vI == nil {
		return true
	}
	if vI.GetPriority() < vJ.GetPriority() {
		return true
	}
	if vI.GetPriority() > vJ.GetPriority() {
		return false
	}
	return vI.GetName() < vJ.GetName()
}

func (instance *Updates) Add(v Update) {
	for i, existing := range *instance {
		if existing.GetName() == v.GetName() {
			(*instance)[i] = v
			sort.Sort(instance)
			return
		}
	}
	*instance = append(*instance, v)
	sort.Sort(instance)
}

type UpdateFunc func(p *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured, argument *string) error

type updateFunc struct {
	UpdateFunc
	transformation
}

func (instance updateFunc) TransformForUpdate(p *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured, argument *string) error {
	return instance.UpdateFunc(p, existing, target, argument)
}

func (instance *Transformations) RegisterUpdate(v Update) error {
	name := v.GetName()
	if _, err := name.MarshalText(); err != nil {
		return err
	}
	if instance == nil {
		*instance = Transformations{}
	}
	for _, other := range instance.Updates {
		if other.GetName() == name {

		}
	}
	instance.Updates.Add(v)
	return nil
}

func (instance *Transformations) MustRegisterUpdate(v Update) {
	if err := instance.RegisterUpdate(v); err != nil {
		panic(err)
	}
}

func (instance *Transformations) RegisterUpdateFunc(name model.TransformationName, v UpdateFunc) error {
	return instance.RegisterUpdate(updateFunc{v, transformation{name}})
}

func (instance *Transformations) MustRegisterUpdateFunc(name model.TransformationName, v UpdateFunc) {
	if err := instance.RegisterUpdateFunc(name, v); err != nil {
		panic(err)
	}
}

func (instance Transformations) TransformForUpdate(p *model.Project, existing unstructured.Unstructured, target *unstructured.Unstructured) error {
	return instance.Updates.TransformForUpdate(p, existing, target)
}
