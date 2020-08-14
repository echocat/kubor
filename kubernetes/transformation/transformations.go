package transformation

import (
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	Default = Transformations{}
)

type Transformations struct {
	Updates Updates
	Creates Creates
}

type Transformation interface {
	GetName() model.TransformationName
	GetPriority() int32
	DefaultEnabled(target *unstructured.Unstructured) bool
}

type transformation struct {
	name model.TransformationName
}

func (instance transformation) GetName() model.TransformationName {
	return instance.name
}

func (instance transformation) DefaultEnabled(*unstructured.Unstructured) bool {
	return true
}

func (instance transformation) GetPriority() int32 {
	return 0
}
