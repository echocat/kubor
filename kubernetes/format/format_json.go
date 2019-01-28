package format

import (
	"io"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

var _ = Provider.MustRegister(VariantJson, &JsonFormat{
	Scheme:      scheme.Scheme,
	MetaFactory: json.DefaultMetaFactory,
})

type JsonFormat struct {
	Scheme      *runtime.Scheme
	MetaFactory json.MetaFactory
}

func (instance JsonFormat) Supports(gvks ...schema.GroupVersionKind) bool {
	for _, gvk := range gvks {
		if !instance.Scheme.Recognizes(gvk) {
			return false
		}
	}
	return true
}

func (instance JsonFormat) Format(to io.Writer, objects ...runtime.Object) error {
	for i, object := range objects {
		if i > 0 {
			if _, err := to.Write([]byte("\n---\n")); err != nil {
				return err
			}
		}
		if err := json.NewSerializer(instance.MetaFactory, instance.Scheme, instance.Scheme, true).
			Encode(object, to); err != nil {
			return err
		}
	}
	return nil
}
