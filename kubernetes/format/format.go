package format

import (
	"errors"
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

type Variant string

const (
	VariantTable = Variant("table")
	VariantYaml  = Variant("yaml")
	VariantJson  = Variant("json")
)

var (
	Provider Registry = &VariantToFormat{}

	ErrUnsupportedVariant = errors.New("variant not supported")
	supportedVariants     = []Variant{VariantTable, VariantYaml, VariantJson}
)

func (instance Variant) String() string {
	return string(instance)
}

func (instance *Variant) Set(plain string) error {
	for _, candidate := range supportedVariants {
		if strings.ToLower(string(candidate)) == strings.ToLower(plain) {
			*instance = candidate
			return nil
		}
	}
	return fmt.Errorf("unsupported format variant: %s", plain)
}

type Format interface {
	Supports(...schema.GroupVersionKind) bool
	Format(to io.Writer, objects ...runtime.Object) error
}

type CombinableFormat interface {
	Format
	ShouldBeCombined() bool
}

type Registry interface {
	Format(variant Variant, to io.Writer, objects ...runtime.Object) error
	Variants() []Variant
	Register(variant Variant, format Format) error
	MustRegister(variant Variant, format Format) Registry
}

type Formats []Format

func (instance Formats) Supports(gvks ...schema.GroupVersionKind) bool {
	for _, gvk := range gvks {
		if !instance.supports(gvk) {
			return false
		}
	}
	return true
}

func (instance Formats) supports(gvk schema.GroupVersionKind) bool {
	for _, format := range instance {
		if format.Supports(gvk) {
			return true
		}
	}
	return false
}

func (instance Formats) Format(to io.Writer, objects ...runtime.Object) error {
	grouped := map[schema.GroupVersionKind][]runtime.Object{}
	for _, object := range objects {
		gvk := NormalizeGroupVersionKind(object.GetObjectKind().GroupVersionKind())
		if group, ok := grouped[gvk]; ok {
			grouped[gvk] = append(group, object)
		} else {
			grouped[gvk] = []runtime.Object{object}
		}
	}
	for gvk, objects := range grouped {
		atLeastOneFound := false
		for _, format := range instance {
			if format.Supports(gvk) {
				atLeastOneFound = true
				if err := format.Format(to, objects...); err != nil {
					return err
				}
				break
			}
		}
		if !atLeastOneFound {
			return fmt.Errorf("there is no formatter that can format %v", FormatGroupVersionKind(gvk))
		}
	}
	return nil
}

type VariantToFormat map[Variant]Format

func (instance VariantToFormat) Format(variant Variant, to io.Writer, objects ...runtime.Object) error {
	if instance == nil {
		return ErrUnsupportedVariant
	}
	if f, ok := instance[variant]; ok {
		return f.Format(to, objects...)
	}
	return ErrUnsupportedVariant
}

func (instance VariantToFormat) Variants() []Variant {
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

func (instance *VariantToFormat) Register(variant Variant, format Format) error {
	if instance == nil {
		*instance = make(VariantToFormat)
	}
	if combinable, ok := format.(CombinableFormat); ok && combinable.ShouldBeCombined() {
		existing := (*instance)[variant]
		if existing == nil {
			existing = Formats{}
		}
		if formats, ok := existing.(Formats); ok {
			(*instance)[variant] = append(formats, format)
			return nil
		} else {
			return fmt.Errorf("more than one time a try to assign a format for variant %v", variant)
		}
	} else {
		if _, existing := (*instance)[variant]; existing {
			return fmt.Errorf("more than one time a try to assign a format for variant %v", variant)
		}
		(*instance)[variant] = format
		return nil
	}
}

func (instance *VariantToFormat) MustRegister(variant Variant, format Format) Registry {
	if err := instance.Register(variant, format); err != nil {
		panic(err)
	} else {
		return instance
	}
}

func NormalizeGroupVersionKind(in schema.GroupVersionKind) schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   strings.ToLower(in.Group),
		Version: strings.ToLower(in.Version),
		Kind:    strings.ToLower(in.Kind),
	}
}

func NormalizeGroupVersionKinds(in []schema.GroupVersionKind) []schema.GroupVersionKind {
	result := make([]schema.GroupVersionKind, len(in))
	for i, val := range in {
		result[i] = NormalizeGroupVersionKind(val)
	}
	return result
}

func FormatGroupVersionKind(in schema.GroupVersionKind) string {
	toFormat := NormalizeGroupVersionKind(in)
	result := toFormat.Kind
	if toFormat.Version != "" {
		result = toFormat.Version + "." + result
	}
	if toFormat.Group != "" {
		result = toFormat.Group + "/" + result
	}
	return result
}
