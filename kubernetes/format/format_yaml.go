package format

import (
	"io"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

var _ = Provider.MustRegister(VariantYaml, &YamlFormat{
	Scheme:      scheme.Scheme,
	MetaFactory: json.DefaultMetaFactory,
})

type YamlFormat struct {
	Scheme      *runtime.Scheme
	MetaFactory json.MetaFactory
}

func (instance YamlFormat) Supports(gvks ...schema.GroupVersionKind) bool {
	for _, gvk := range gvks {
		if !instance.Scheme.Recognizes(gvk) {
			return false
		}
	}
	return true
}

func (instance YamlFormat) Format(to io.Writer, objects ...runtime.Object) error {
	for i, object := range objects {
		if i > 0 {
			if _, err := to.Write([]byte("\n---\n")); err != nil {
				return err
			}
		}
		if err := json.NewYAMLSerializer(instance.MetaFactory, instance.Scheme, instance.Scheme).
			Encode(object, to); err != nil {
			return err
		}
	}
	return nil
}
