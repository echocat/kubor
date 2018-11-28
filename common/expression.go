package common

import (
	"fmt"
	"reflect"
	"strings"
)

func MustEvaluateExpression(expression string, object interface{}) interface{} {
	result, err := EvaluateExpression(expression, object)
	if err != nil {
		panic(err)
	}
	return result
}

func EvaluateExpression(expression string, object interface{}) (interface{}, error) {
	trimmedExpression := strings.TrimSpace(expression)
	if trimmedExpression == "" {
		return nil, fmt.Errorf("empty expression")
	}
	v := reflect.ValueOf(object)
	parts := strings.Split(expression, ".")
	for i := 0; i < len(parts); i++ {
		part := strings.TrimSpace(parts[i])
		if part == "" {
			return nil, fmt.Errorf("empty part")
		}
		newV, err := getPropertyForFieldName(v, part)
		if err != nil {
			return nil, err
		}
		kind := newV.Kind()
		if kind == reflect.Struct || kind == reflect.Map {
			if i+1 < len(parts) {
				v = newV
			} else if newV.IsValid() {
				return newV.Interface(), nil
			} else {
				return nil, nil
			}
		} else if i+1 == len(parts) && newV.IsValid() {
			return newV.Interface(), nil
		} else {
			return nil, nil
		}
	}
	return nil, nil
}

func getPropertyForFieldName(v reflect.Value, name string) (reflect.Value, error) {
	t := v.Type()
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return reflect.Value{}, nil
		}
		v = v.Elem()
	}
	if !v.IsValid() {
		return reflect.Value{}, nil
	}
	for v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		if field, ok := t.FieldByName(name); ok {
			fieldValue := v.FieldByIndex(field.Index)
			return simplifyValue(fieldValue), nil
		}
		return reflect.Value{}, nil
	case reflect.Map:
		value := v.MapIndex(reflect.ValueOf(name))
		if !value.IsValid() {
			return reflect.Value{}, nil
		}
		return simplifyValue(value), nil
	}
	return reflect.Value{}, fmt.Errorf("cannot handle value of type %v", t)
}

func simplifyValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	if !v.IsValid() {
		return reflect.Value{}
	}
	for v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if !v.IsValid() {
		return reflect.Value{}
	}
	return v
}
