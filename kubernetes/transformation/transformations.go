package transformation

import (
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	UpdatesTransformations []UpdateTransformation
	CreateTransformations  []CreateTransformation
)

type UpdateTransformation interface {
	TransformForUpdate(p *model.Project, original unstructured.Unstructured, target *unstructured.Unstructured) error
}

type CreateTransformation interface {
	TransformForCreate(p *model.Project, target *unstructured.Unstructured) error
}

type UpdateTransformationFunc func(p *model.Project, original unstructured.Unstructured, target *unstructured.Unstructured) error

func (instance UpdateTransformationFunc) TransformForUpdate(p *model.Project, original unstructured.Unstructured, target *unstructured.Unstructured) error {
	return instance(p, original, target)
}

type CreateTransformationFunc func(p *model.Project, target *unstructured.Unstructured) error

func (instance CreateTransformationFunc) TransformForCreate(p *model.Project, target *unstructured.Unstructured) error {
	return instance(p, target)
}

func RegisterUpdateTransformation(v UpdateTransformation) UpdateTransformation {
	UpdatesTransformations = append(UpdatesTransformations, v)
	return v
}

func RegisterUpdateTransformationFunc(v UpdateTransformationFunc) {
	RegisterUpdateTransformation(v)
}

func RegisterCreateTransformation(v CreateTransformation) CreateTransformation {
	CreateTransformations = append(CreateTransformations, v)
	return v
}

func RegisterCreateTransformationFunc(v CreateTransformationFunc) {
	RegisterCreateTransformation(v)
}

func TransformForUpdate(p *model.Project, original unstructured.Unstructured, target *unstructured.Unstructured) error {
	for _, transformation := range UpdatesTransformations {
		if err := transformation.TransformForUpdate(p, original, target); err != nil {
			return err
		}
	}
	return nil
}

func TransformForCreate(p *model.Project, target *unstructured.Unstructured) error {
	for _, transformation := range CreateTransformations {
		if err := transformation.TransformForCreate(p, target); err != nil {
			return err
		}
	}
	return nil
}
