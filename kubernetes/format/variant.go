package format

import (
	"errors"
	"fmt"
	"strings"
)

type Variant string

const (
	VariantTable = Variant("table")
	VariantYaml  = Variant("yaml")
	VariantJson  = Variant("json")
)

var (
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
