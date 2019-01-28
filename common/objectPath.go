package common

import (
	"reflect"
)

func GetObjectPathValue(object interface{}, segments ...string) interface{} {
	if len(segments) == 0 {
		return object
	}
	v := reflect.ValueOf(object)
	for i := 0; i < len(segments); i++ {
		segment := segments[i]
		newV := getPropertyForFieldName(v, segment)
		kind := newV.Kind()
		if kind == reflect.Struct || kind == reflect.Map {
			if i+1 < len(segments) {
				v = newV
			} else if newV.IsValid() {
				return newV.Interface()
			} else {
				return nil
			}
		} else if i+1 == len(segments) && newV.IsValid() {
			return newV.Interface()
		} else {
			return nil
		}
	}
	return nil
}

func getPropertyForFieldName(v reflect.Value, name string) reflect.Value {
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
	t := v.Type()
	switch v.Kind() {
	case reflect.Struct:
		if field, ok := t.FieldByName(name); ok {
			fieldValue := v.FieldByIndex(field.Index)
			return simplifyValue(fieldValue)
		}
		return reflect.Value{}
	case reflect.Map:
		value := v.MapIndex(reflect.ValueOf(name))
		if !value.IsValid() {
			return reflect.Value{}
		}
		return simplifyValue(value)
	}
	return reflect.Value{}
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
