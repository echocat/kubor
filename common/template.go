package common

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"reflect"
	"text/template"
)

func NewTemplate(name string, plain string) (*template.Template, error) {
	return template.
		New(name).
		Funcs(sprig.HermeticTxtFuncMap()).
		Funcs(template.FuncMap{
			"optional": templateFuncOptional,
			"opt":      templateFuncOptional,
		}).
		Option("missingkey=error").
		Parse(plain)
}

func templateFuncDefault(name string, def interface{}, holder interface{}) (interface{}, error) {
	v := reflect.ValueOf(holder)
	t := v.Type()
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.IsNil() {
		return def, nil
	}
	for v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		if field, ok := t.FieldByName(name); ok {
			fieldValue := v.FieldByIndex(field.Index)
			return fieldValue.Interface(), nil
		}
		return def, nil
	case reflect.Map:
		value := v.MapIndex(reflect.ValueOf(name))
		if !value.IsValid() {
			return nil, nil
		}
		return value.Interface(), nil
	}
	return def, fmt.Errorf("cannot get value for '%s' because cannot handle values of type %v", name, t)
}

func templateFuncOptional(name string, holder interface{}) (interface{}, error) {
	return templateFuncDefault(name, nil, holder)
}
