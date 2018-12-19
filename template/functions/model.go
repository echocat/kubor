package functions

import (
	"bytes"
	"fmt"
	"kubor/template"
	"reflect"
	"strings"
)

var (
	errType              = reflect.TypeOf((*error)(nil)).Elem()
	executionContextType = reflect.TypeOf((*template.ExecutionContext)(nil)).Elem()
)

type Function struct {
	Func        interface{} `yaml:"-" json:"-"`
	Name        string      `yaml:"name" json:"name"`
	Category    string      `yaml:"category" json:"category"`
	Description string      `yaml:"description,omitempty" json:"description,omitempty"`
	Parameters  Parameters  `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Returns     Return      `yaml:"returns" json:"returns"`
}

func (instance Function) Complete() (result Function, err error) {
	ft := reflect.TypeOf(instance.Func)
	if ft.Kind() != reflect.Func {
		return Function{}, fmt.Errorf("provided Func field is not of expected type func - it is: %v", ft)
	}
	result = instance
	if result.Parameters, err = instance.completeParameters(ft); err != nil {
		return
	}
	if result.Returns, err = instance.completeReturns(ft); err != nil {
		return
	}
	return
}

func (instance Function) completeReturns(ft reflect.Type) (Return, error) {
	if ft.NumOut() <= 0 {
		return Return{}, fmt.Errorf("function does not have any returning arguments")
	}
	if ft.NumOut() > 2 {
		return Return{}, fmt.Errorf("function does have too many returning argument: %d", ft.NumOut())
	}
	if ft.NumOut() == 2 && ft.Out(1) != errType {
		return Return{}, fmt.Errorf("2nd argument of function has to be of type error but is: %v", ft.Out(1))
	}
	result := instance.Returns
	result.Type = NormalizeType(ft.Out(0))

	return result, nil
}

func (instance Function) completeParameters(ft reflect.Type) (Parameters, error) {
	result := instance.Parameters

	var actualIndex int
	for i := 0; i < ft.NumIn(); i++ {
		pt := ft.In(i)
		if pt != executionContextType {
			if len(result) <= actualIndex {
				return Parameters{}, fmt.Errorf("too few number of parameters defined - #%d missing", actualIndex)
			}
			existing := result[actualIndex]
			if parameter, err := instance.completeParameter(ft, i, existing); err != nil {
				return Parameters{}, err
			} else {
				result[actualIndex] = parameter
			}
			actualIndex++
		}
	}

	return result, nil
}

func (instance Function) completeParameter(ft reflect.Type, index int, existing Parameter) (Parameter, error) {
	pt := ft.In(index)

	result := existing
	result.Type = NormalizeType(pt)
	result.VarArg = ft.IsVariadic() && ft.NumIn() == index+1

	return result, nil
}

func (instance Function) String() string {
	return fmt.Sprintf("%s(%s): %s",
		instance.GetName(),
		instance.GetParameters().String(),
		instance.GetReturns().String(),
	)
}

func (instance Function) GetName() string {
	return instance.Name
}

func (instance Function) GetCategory() string {
	if instance.Category == "" {
		return "general"
	}
	return instance.Category
}

func (instance Function) GetDescription() string {
	return instance.Description
}

func (instance Function) GetParameters() Parameters {
	return instance.Parameters
}

func (instance Function) GetReturns() Return {
	return instance.Returns
}

type Parameter struct {
	Name        string `yaml:"name" json:"name"`
	Type        string `yaml:"type" json:"type"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	VarArg      bool   `yaml:"varArg,omitempty" json:"varArg,omitempty"`
}

func (instance Parameter) String() string {
	return fmt.Sprintf("%s:%s", instance.GetName(), instance.GetType())
}

func (instance Parameter) GetType() string {
	t := instance.Type
	if strings.HasPrefix(t, "[]") && instance.IsVarArg() {
		t = "..." + t[2:]
	}
	return t
}

func (instance Parameter) GetName() string {
	return instance.Name
}

func (instance Parameter) GetDescription() string {
	return instance.Description
}

func (instance Parameter) IsVarArg() bool {
	return instance.VarArg
}

type Return struct {
	Type        string `yaml:"type" json:"type"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

func (instance Return) String() string {
	return instance.GetType()
}

func (instance Return) GetType() string {
	return instance.Type
}

func (instance Return) GetDescription() string {
	return instance.Description
}

type Parameters []Parameter

func (instance Parameters) String() string {
	buf := new(bytes.Buffer)
	for i, parameter := range instance {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(parameter.String())
	}
	return buf.String()
}
