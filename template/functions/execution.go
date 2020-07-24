package functions

import (
	"fmt"
	"github.com/echocat/kubor/template"
	"reflect"
)

func (instance Function) Execute(context template.ExecutionContext, args ...interface{}) (interface{}, error) {
	fv := reflect.ValueOf(instance.f)
	ft := fv.Type()

	if pvs, err := instance.createExecutionArguments(ft, context, args...); err != nil {
		return nil, err
	} else {
		rvs := instance.call(fv, pvs...)
		return instance.evaluateResult(rvs...)
	}
}

func (instance Function) call(fv reflect.Value, in ...reflect.Value) []reflect.Value {
	if fv.Type().IsVariadic() {
		return fv.CallSlice(in)
	}
	return fv.Call(in)
}

func (instance Function) evaluateResult(out ...reflect.Value) (interface{}, error) {
	if len(out) <= 0 {
		return Return{}, fmt.Errorf("function does not have any returning arguments")
	}
	if len(out) > 2 {
		return Return{}, fmt.Errorf("function does have too many returning argument: %d", len(out))
	}
	if len(out) == 2 && !out[1].IsNil() {
		if err, ok := out[1].Interface().(error); ok {
			return nil, err
		}
		return nil, fmt.Errorf("2nd argument of function has to be of type error but is: %v", out[1])
	}
	return out[0].Interface(), nil
}

func (instance Function) createExecutionArguments(ft reflect.Type, context template.ExecutionContext, args ...interface{}) ([]reflect.Value, error) {
	result := make([]reflect.Value, ft.NumIn())

	variadic := false
	numberOfRequiredParameters := 0
	for i := 0; i < ft.NumIn(); i++ {
		switch ft.In(i) {
		case executionContextType:
		default:
			if !ft.IsVariadic() {
				numberOfRequiredParameters++
			} else {
				variadic = true
			}
		}
	}

	if variadic {
		if len(args) < numberOfRequiredParameters {
			return []reflect.Value{}, fmt.Errorf("wrong number of args: need at least %d got %d", numberOfRequiredParameters, len(args))
		}
	} else if len(args) != numberOfRequiredParameters {
		return []reflect.Value{}, fmt.Errorf("wrong number of args: want %d got %d", numberOfRequiredParameters, len(args))
	}

	var argIndex int
	for i := 0; i < ft.NumIn(); i++ {
		pt := ft.In(i)
		if pt == executionContextType {
			result[i] = reflect.ValueOf(context)
		} else if ft.IsVariadic() {
			if pv, err := instance.createExecutionVarargArgument(argIndex, pt, args[argIndex:]); err != nil {
				return []reflect.Value{}, err
			} else {
				result[i] = pv
				break
			}
		} else {
			if pv, err := instance.createExecutionArgument(argIndex, pt, args[argIndex]); err != nil {
				return []reflect.Value{}, err
			} else {
				result[i] = pv
				argIndex++
			}
		}
	}

	return result, nil
}

func (instance Function) createExecutionArgument(index int, pt reflect.Type, arg interface{}) (reflect.Value, error) {
	valOf := func(pt reflect.Type, in interface{}) reflect.Value {
		if in != nil {
			return reflect.ValueOf(in)
		}
		return reflect.New(pt).Elem()
	}
	av := valOf(pt, arg)
	at := av.Type()
	if !at.AssignableTo(pt) {
		return reflect.Value{}, fmt.Errorf("%v is not assignable to %v for argument #%d", at, pt, index)
	}
	return av, nil
}

func (instance Function) createExecutionVarargArgument(index int, pt reflect.Type, args []interface{}) (reflect.Value, error) {
	valOf := func(pt reflect.Type, in interface{}) reflect.Value {
		if in != nil {
			return reflect.ValueOf(in)
		}
		return reflect.New(pt).Elem()
	}
	av := reflect.MakeSlice(pt, len(args), len(args))
	for i := 0; i < len(args); i++ {
		av.Index(i).Set(valOf(pt, args[i]))
	}
	at := av.Type()
	if !at.AssignableTo(pt) {
		return reflect.Value{}, fmt.Errorf("%v is not assignable to %v for argument #%d", at, pt, index)
	}
	return av, nil
}
