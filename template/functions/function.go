package functions

import (
	"bytes"
	"fmt"
	"kubor/template"
	"reflect"
	"regexp"
	"strings"
)

var (
	errType              = reflect.TypeOf((*error)(nil)).Elem()
	executionContextType = reflect.TypeOf((*template.ExecutionContext)(nil)).Elem()
)

type Function struct {
	Description string     `yaml:"description,omitempty" json:"description,omitempty"`
	Parameters  Parameters `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Returns     Return     `yaml:"returns" json:"returns"`

	f interface{}
}

func (instance Function) MustWithFunc(f interface{}) Function {
	if result, err := instance.WithFunc(f); err != nil {
		panic(err)
	} else {
		return result
	}
}

func (instance Function) WithFunc(f interface{}) (result Function, err error) {
	ft := reflect.TypeOf(f)
	if ft.Kind() != reflect.Func {
		return Function{}, fmt.Errorf("provided value is not of expected type func - it is: %v", ft)
	}
	result = instance
	result.f = f
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
	return fmt.Sprintf("func(%s): %s",
		instance.GetParameters().String(),
		instance.GetReturns().String(),
	)
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

func (instance Function) MatchesFulltextSearch(term *regexp.Regexp) bool {
	if term.FindStringIndex(instance.Description) != nil {
		return true
	}
	if instance.Parameters.MatchesFulltextSearch(term) {
		return true
	}
	if instance.Returns.MatchesFulltextSearch(term) {
		return true
	}
	return false
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

func (instance Parameter) MatchesFulltextSearch(term *regexp.Regexp) bool {
	if term.FindStringIndex(instance.Name) != nil {
		return true
	}
	if term.FindStringIndex(instance.Description) != nil {
		return true
	}
	return false
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

func (instance Return) MatchesFulltextSearch(term *regexp.Regexp) bool {
	if term.FindStringIndex(instance.Description) != nil {
		return true
	}
	return false
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

func (instance Parameters) MatchesFulltextSearch(term *regexp.Regexp) bool {
	for _, parameter := range instance {
		if parameter.MatchesFulltextSearch(term) {
			return true
		}
	}
	return false
}

type Functions map[string]Function

func (instance Functions) GetFunctions() (template.Functions, error) {
	result := template.Functions{}
	for name, function := range instance {
		result[name] = function
	}
	return result, nil
}

func (instance Functions) FilterBy(predicate FunctionPredicate) (Functions, error) {
	result := Functions{}
	for name, candidate := range instance {
		if predicate == nil {
			result[name] = candidate
		} else if ok, err := predicate(candidate); err != nil {
			return Functions{}, err
		} else if ok {
			result[name] = candidate
		}
	}
	return result, nil
}

func (instance Functions) With(functionName string, function Function) Functions {
	result := instance
	if result == nil {
		result = Functions{}
	}
	result[functionName] = function
	return result
}

func (instance Functions) Without(functionName string) Functions {
	result := instance
	if result == nil {
		result = Functions{}
	}
	delete(result, functionName)
	return result
}

type FunctionPredicate func(candidate Function) (bool, error)
