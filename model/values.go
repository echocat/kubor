package model

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"kubor/common"
	"os"
	"strings"
)

type Values map[string]interface{}

func (instance Values) MergeWith(input ...Values) Values {
	result := instance
	if result == nil {
		result = Values{}
	}
	for _, values := range input {
		for key, value := range values {
			result[key] = value
		}
	}
	return result
}

func (instance Values) MergeWithFiles(files ...string) (Values, error) {
	result := instance
	for _, file := range files {
		if newResult, err := result.MergeWithFile(file); err != nil {
			return nil, err
		} else {
			result = newResult
		}
	}
	return result, nil
}

func (instance Values) MergeWithFile(file string) (Values, error) {
	if f, err := os.Open(file); err != nil {
		return nil, fmt.Errorf("cannot merge with '%s': %v", file, err)
	} else {
		//noinspection GoUnhandledErrorResult
		defer f.Close()
		tmp := Values{}
		if err := yaml.NewDecoder(f).Decode(&tmp); err != nil {
			return nil, fmt.Errorf("cannot merge with '%s': %v", file, err)
		} else {
			return instance.MergeWith(tmp), nil
		}
	}
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

type ConditionalValues struct {
	On     common.EvaluatingPredicate `yaml:"on,omitempty" json:"on,omitempty"`
	Values Values                     `yaml:",inline" json:",inline"`
}
