package transformation

import (
	"fmt"
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sort"
)

type Create interface {
	Transformation
	TransformForCreate(p *model.Project, target *unstructured.Unstructured, argument *string) error
}

type Creates []Create

func (instance Creates) TransformForCreate(p *model.Project, target *unstructured.Unstructured) error {
	for _, transformation := range instance {
		name := transformation.GetName()
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
	return nil
}

func (instance Creates) Len() int      { return len(instance) }
func (instance Creates) Swap(i, j int) { instance[i], instance[j] = instance[j], instance[i] }
func (instance Creates) Less(i, j int) bool {
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

func (instance *Creates) Add(v Create) {
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

type CreateFunc func(p *model.Project, target *unstructured.Unstructured, argument *string) error

type createFunc struct {
	CreateFunc
	transformation
}

func (instance createFunc) TransformForCreate(p *model.Project, target *unstructured.Unstructured, argument *string) error {
	return instance.CreateFunc(p, target, argument)
}

func (instance *Transformations) RegisterCreate(v Create) error {
	name := v.GetName()
	if _, err := name.MarshalText(); err != nil {
		return err
	}
	if instance == nil {
		*instance = Transformations{}
	}
	instance.Creates.Add(v)
	return nil
}

func (instance *Transformations) MustRegisterCreate(v Create) {
	if err := instance.RegisterCreate(v); err != nil {
		panic(err)
	}
}

func (instance *Transformations) RegisterCreateFunc(name model.TransformationName, v CreateFunc) error {
	return instance.RegisterCreate(createFunc{v, transformation{name}})
}

func (instance *Transformations) MustRegisterCreateFunc(name model.TransformationName, v CreateFunc) {
	if err := instance.RegisterCreateFunc(name, v); err != nil {
		panic(err)
	}
}

func (instance Transformations) TransformForCreate(p *model.Project, target *unstructured.Unstructured) error {
	return instance.Creates.TransformForCreate(p, target)
}
