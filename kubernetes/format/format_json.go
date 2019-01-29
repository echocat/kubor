package format

import (
	"io"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

var _ = MustRegister(VariantJson, NewJsonFormat())

func NewJsonFormat() *JsonFormat {
	return &JsonFormat{
		Encoder: json.NewSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, true),
	}
}

type JsonFormat struct {
	Encoder runtime.Encoder
}

func (instance JsonFormat) Format(to io.Writer, supplier ObjectSupplier) error {
	first := true
	for {
		object, err := supplier()
		if err != nil || object == nil {
			return err
		}
		if first {
			first = false
		} else if _, err := to.Write([]byte("\n---\n")); err != nil {
			return err
		}
		if err := instance.Encoder.Encode(object, to); err != nil {
			return err
		}
	}
}
