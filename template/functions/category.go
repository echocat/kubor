package functions

import (
	"github.com/echocat/kubor/template"
)

type Category struct {
	Functions Functions
}

func (instance Category) GetFunctions() (template.Functions, error) {
	return instance.Functions.GetFunctions()
}

func (instance Category) FilterBy(predicate FunctionPredicate) (Category, error) {
	if functions, err := instance.Functions.FilterBy(predicate); err != nil {
		return Category{}, err
	} else {
		return Category{
			Functions: functions,
		}, nil
	}
}

func (instance Category) With(functionName string, function Function) Category {
	result := instance
	result.Functions = result.Functions.With(functionName, function)
	return result
}

func (instance Category) Without(functionName string) Category {
	result := instance
	result.Functions = result.Functions.Without(functionName)
	return result
}

type Categories map[string]Category

func (instance Categories) GetFunctions() (template.Functions, error) {
	result := template.Functions{}
	for _, category := range instance {
		if functions, err := category.GetFunctions(); err != nil {
			return template.Functions{}, err
		} else {
			for name, function := range functions {
				result[name] = function
			}
		}
	}
	return result, nil
}

func (instance Categories) With(categoryName string, category Category) Categories {
	result := instance
	result[categoryName] = category
	return result
}

func (instance Categories) Without(categoryName string) Categories {
	result := instance
	if result == nil {
		result = Categories{}
	}
	delete(result, categoryName)
	return result
}

func (instance Categories) WithFunction(categoryName string, functionName string, function Function) Categories {
	result := instance
	if result == nil {
		result = Categories{}
	}
	result[categoryName] = result[categoryName].
		With(functionName, function)
	return result
}

func (instance Categories) WithoutFunction(functionName string) Categories {
	result := instance
	if result == nil {
		result = Categories{}
	}
	for name, category := range result {
		newCategory := category.Without(functionName)
		if len(newCategory.Functions) > 0 {
			result[name] = newCategory
		} else {
			delete(result, name)
		}
	}
	return result
}

func (instance Categories) FilterBy(predicate FunctionPredicate) (Categories, error) {
	result := Categories{}
	for name, category := range instance {
		if funcs, err := category.FilterBy(predicate); err != nil {
			return Categories{}, err
		} else if len(funcs.Functions) > 0 {
			result[name] = funcs
		}
	}
	return result, nil
}

var CategoriesDefault = Categories{
	"codecs":        CategoryCodecs,
	"conversations": CategoryConversations,
	"general":       CategoryGeneral,
	"kubernetes":    CategoryKubernetes,
	"math":          CategoryMath,
	"path":          CategoryPath,
	"regexp":        CategoryRegexp,
	"serialization": CategorySerialization,
	"strings":       CategoryStrings,
	"templating":    CategoryTemplating,
}
