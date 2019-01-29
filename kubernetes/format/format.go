package format

import (
	"fmt"
	"github.com/levertonai/kubor/kubernetes"
	"io"
)

var (
	DefaultFormats Formats = &formats{}
)

type ObjectSupplier func() (kubernetes.Object, error)

type Format interface {
	Format(to io.Writer) (Task, error)
}

type Task interface {
	io.Closer
	Next(kubernetes.Object) error
}

type Formats interface {
	SupportedVariants() []Variant
	Get(variant Variant) (Format, error)
	Register(variant Variant, format Format) error
}

func Get(variant Variant) (Format, error) {
	return DefaultFormats.Get(variant)
}

func Register(variant Variant, format Format) error {
	return DefaultFormats.Register(variant, format)
}

func MustRegister(variant Variant, format Format) Formats {
	if err := Register(variant, format); err != nil {
		panic(err)
	} else {
		return DefaultFormats
	}
}

type formats map[Variant]Format

func (instance formats) SupportedVariants() []Variant {
	if instance == nil {
		return []Variant{}
	}
	result := make([]Variant, len(instance))
	var i int
	for variant := range instance {
		result[i] = variant
		i++
	}
	return result
}

func (instance formats) Get(variant Variant) (Format, error) {
	if instance == nil {
		return nil, ErrUnsupportedVariant
	}
	if format, ok := instance[variant]; ok {
		return format, nil
	}
	return nil, ErrUnsupportedVariant
}

func (instance *formats) Register(variant Variant, format Format) error {
	if instance == nil {
		*instance = make(formats)
	}
	if _, exists := (*instance)[variant]; exists {
		return fmt.Errorf("more than one time a try to assign a format for variant %v", variant)
	}
	(*instance)[variant] = format
	return nil
}
