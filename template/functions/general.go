package functions

import (
	"fmt"
	"reflect"
	"strings"
)

var _ = Register(Function{
	Name:        "optional",
	Category:    "general",
	Description: "Asks the given <holder> if a property of given <name> exists and returns it. Otherwise the result is <nil>.",
	Parameters: Parameters{{
		Name: "name",
	}, {
		Name: "holder",
	}},
	Func: func(name string, holder interface{}) (interface{}, error) {
		v := reflect.ValueOf(holder)
		t := v.Type()
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.IsNil() {
			return nil, nil
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
			return nil, nil
		case reflect.Map:
			value := v.MapIndex(reflect.ValueOf(name))
			if !value.IsValid() {
				return nil, nil
			}
			return value.Interface(), nil
		}
		return nil, fmt.Errorf("cannot get value for '%s' because cannot handle values of type %v", name, t)
	},
}, Function{
	Name:        "contains",
	Category:    "general",
	Description: `Checks if the given <input> string contains the given <search> string.`,
	Parameters: Parameters{{
		Name: "search",
	}, {
		Name: "input",
	}},
	Func: func(search interface{}, input interface{}) (bool, error) {
		if s, ok := input.(string); ok {
			return strings.Contains(s, fmt.Sprint(search)), nil
		} else if s, ok := input.(*string); ok {
			return strings.Contains(*s, fmt.Sprint(search)), nil
		}
		return false, fmt.Errorf("currently contains only supports strings but got: %v", reflect.TypeOf(input))
	},
})
