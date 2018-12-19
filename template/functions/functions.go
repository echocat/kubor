package functions

import (
	"fmt"
	"kubor/template"
	"strings"
	"sync"
)

var globalRegistry = &Functions{}

func GlobalRegistry() *Functions {
	return globalRegistry
}

type Functions struct {
	functions map[string]*Function

	mutex sync.Mutex
}

func (instance *Functions) Has(name string) bool {
	functions := instance.functions
	if functions == nil {
		return false
	}
	return functions[name] != nil
}

func (instance *Functions) Get(name string) *Function {
	functions := instance.functions
	if functions == nil {
		return nil
	}
	if result := functions[name]; result != nil {
		f := *result
		return &f
	}
	return nil
}

func (instance *Functions) Add(functions ...Function) (*Functions, error) {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()

	if instance.functions == nil {
		instance.functions = map[string]*Function{}
	}

	for _, function := range functions {
		if completed, err := function.Complete(); err != nil {
			return nil, fmt.Errorf("cannot add function %s: %v", function.GetName(), err)
		} else {
			instance.functions[function.Name] = &completed
		}
	}

	return instance, nil
}

func (instance *Functions) Remove(names ...string) *Functions {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()

	if instance.functions == nil {
		return instance
	}

	for _, name := range names {
		delete(instance.functions, name)
	}

	return instance
}

func (instance *Functions) GetAll() []Function {
	return instance.GetAllBy(nil)
}

func (instance *Functions) GetAllOfCategory(category string) []Function {
	expectedCategory := strings.ToLower(category)

	return instance.GetAllBy(func(candidate Function) bool {
		return strings.ToLower(candidate.GetCategory()) == expectedCategory
	})
}

func (instance *Functions) GetAllBy(predicate func(Function) bool) (result []Function) {
	functions := instance.functions
	if functions == nil {
		return
	}

	for _, function := range functions {
		uFunction := *function
		if predicate == nil || predicate(uFunction) {
			result = append(result, uFunction)
		}
	}

	return
}

func (instance *Functions) GetFunctions() (result template.Functions, err error) {
	functions := instance.functions
	if functions == nil {
		return
	}

	result = make(template.Functions, len(functions))
	var i int
	for _, function := range functions {
		result[i] = *function
		i++
	}

	return
}

func Register(functions ...Function) *Functions {
	if result, err := GlobalRegistry().Add(functions...); err != nil {
		panic(err)
	} else {
		return result
	}
}
