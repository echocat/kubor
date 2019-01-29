package format

import (
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

var _ = MustRegister(VariantYaml, NewYamlFormat())

func NewYamlFormat() *JsonFormat {
	return &JsonFormat{
		Encoder: json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme),
	}
}
