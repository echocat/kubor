package functions

import (
	"github.com/echocat/kubor/template"
	"reflect"
	"strings"
)

func NormalizeType(t reflect.Type) string {
	result := t.String()
	result = strings.Replace(result, " ", "", -1)
	result = strings.Replace(result, "interface{}", "any", -1)
	return result
}

func DefaultTemplateFactory() template.Factory {
	return &template.FactoryImpl{
		FunctionProvider: CategoriesDefault,
	}
}
