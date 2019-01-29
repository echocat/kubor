package functions

import (
	"fmt"
	"github.com/levertonai/kubor/template"
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

	numberOfParameters := 0
	for i := 0; i < ft.NumIn(); i++ {
		switch ft.In(i) {
		case executionContextType:
		default:
			numberOfParameters++
		}
	}

	if len(args) != numberOfParameters {
		return []reflect.Value{}, fmt.Errorf("wrong number of args: want %d got %d", numberOfParameters, len(args))
	}

	var argIndex int
	for i := 0; i < ft.NumIn(); i++ {
		pt := ft.In(i)
		if pt == executionContextType {
			result[i] = reflect.ValueOf(context)
		} else if pv, err := instance.createExecutionArgument(argIndex, ft, pt, args[argIndex]); err != nil {
			return []reflect.Value{}, err
		} else {
			result[i] = pv
			argIndex++
		}
	}

	return result, nil
}

func (instance Function) createExecutionArgument(index int, ft reflect.Type, pt reflect.Type, arg interface{}) (reflect.Value, error) {
	if arg == nil {
		return reflect.New(pt), nil
	}
	av := reflect.ValueOf(arg)
	at := av.Type()
	if !at.AssignableTo(pt) {
		return reflect.Value{}, fmt.Errorf("%v is not assignable to %v for argument #%d", at, pt, index)
	}
	return av, nil
}
