package model

import (
	"fmt"
	"kubor/common"
	"strings"
)

type Values map[string]interface{}

func (instance Values) MergeWith(input ...Values) Values {
	result := Values{}
	if instance != nil {
		for k, v := range instance {
			result[k] = v
		}
	}
	for _, values := range input {
		for key, value := range values {
			result[key] = value
		}
	}
	return result
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
