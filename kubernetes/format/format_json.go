package format

import (
	"github.com/levertonai/kubor/kubernetes"
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

func (instance JsonFormat) Format(to io.Writer) (Task, error) {
	return &jsonTask{
		JsonFormat: instance,
		to:         to,
		first:      true,
	}, nil
}

type jsonTask struct {
	JsonFormat
	to    io.Writer
	first bool
}

func (instance *jsonTask) Next(object kubernetes.Object) error {
	if instance.first {
		instance.first = false
	} else if _, err := instance.to.Write([]byte("\n---\n")); err != nil {
		return err
	}
	return instance.Encoder.Encode(object, instance.to)
}

func (instance *jsonTask) Close() error {
	return nil
}
