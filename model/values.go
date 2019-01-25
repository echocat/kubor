package model

import (
	"errors"
	"fmt"
	"github.com/levertonai/kubor/common"
	"gopkg.in/yaml.v2"
	"reflect"
	"strings"
)

type Values map[string]interface{}

func (instance Values) Clone() Values {
	result := Values{}
	if instance != nil {
		for k, v := range instance {
			result[k] = v
		}
	}
	return result
}

func (instance Values) MergeWith(input ...Values) Values {
	result := instance.Clone()
	for _, values := range input {
		for key, value := range values {
			result[key] = value
		}
	}
	return result
}

func (instance *Values) IsCumulative() bool {
	return true
}

func (instance *Values) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	if *instance == nil {
		*instance = Values{}
	}
	if len(parts) > 1 {
		(*instance)[parts[0]] = parts[1]
	} else {
		(*instance)[parts[0]] = ""
	}
	return nil
}

// String returns a readable representation of this value (for usage defaults)
func (instance *Values) String() string {
	return fmt.Sprintf("%s", *instance)
}

// Get returns the slice of strings set by this flag
func (instance *Values) Get() interface{} {
	return *instance
}

type ValuesDefinition struct {
	on     common.EvaluatingPredicate
	values []ValuesDefinitionEntry
}

type ValuesDefinitionEntry struct {
	Key   string
	Value interface{}
}

func (instance *ValuesDefinition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var sym yaml.MapSlice
	if err := unmarshal(&sym); err != nil {
		return err
	}
	result := ValuesDefinition{
		values: []ValuesDefinitionEntry{},
	}
	for _, entry := range sym {
		key := fmt.Sprint(entry.Key)
		if key == "on" || entry.Key == true {
			if err := result.on.Set(fmt.Sprint(entry.Value)); err != nil {
				return err
			}
		} else {
			switch v := entry.Value.(type) {
			case string, uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64, bool:
				result.values = append(result.values, ValuesDefinitionEntry{
					Key:   key,
					Value: fmt.Sprint(v),
				})
			case yaml.MapSlice:
				return errors.New("values with sub-objects are not supported")
			default:
				return fmt.Errorf("values of type %v are not supported", reflect.TypeOf(v))
			}
		}
	}
	*instance = result
	return nil
}

func (instance ValuesDefinition) MarshalYAML() (interface{}, error) {
	result := make(yaml.MapSlice, len(instance.values))
	for i, entry := range instance.values {
		result[i] = yaml.MapItem{
			Key:   entry.Key,
			Value: entry.Value,
		}
	}
	onStr := instance.on.String()
	if onStr != "" {
		result = append(result, yaml.MapItem{
			Key:   "on",
			Value: onStr,
		})
	}
	return result, nil
}

func (instance ValuesDefinition) eval(project Project, override ...Values) (Values, error) {
	result := Values{}

	for _, keyAndValue := range instance.values {
		jobProject := project
		jobProject.Values = jobProject.Values.MergeWith(result)
		jobProject.Values = jobProject.Values.MergeWith(override...)
		key := keyAndValue.Key
		value := keyAndValue.Value

		if pstr, ok := value.(*string); ok {
			value = *pstr
		}

		if str, ok := value.(string); ok {
			if content, err := evaluateTemplate("value."+key, str, jobProject); err != nil {
				return Values{}, err
			} else {
				value = content
			}
		}

		result[key] = value
	}

	return result.MergeWith(override...), nil
}

type ValuesDefinitions []ValuesDefinition

func (instance ValuesDefinitions) ToValues(project Project, override Values) (result Values, err error) {
	for _, definition := range instance {
		if matches, err := definition.on.Matches(project); err != nil {
			return Values{}, err
		} else if matches {
			if result, err = definition.eval(project, result, override); err != nil {
				return Values{}, err
			}
		}
	}
	return
}
