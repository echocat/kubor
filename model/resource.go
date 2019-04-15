package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unicode"
)

type ResourceIdentifier string

func (instance ResourceIdentifier) String() string {
	return string(instance)
}

func (instance *ResourceIdentifier) Set(plain string) error {
	for _, ch := range plain {
		if !((ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-') {
			return fmt.Errorf("illegal %s: character %c is not supported", instance.TypeName(), ch)
		}
	}
	*instance = ResourceIdentifier(plain)
	return nil
}

func (instance ResourceIdentifier) MarshalText() (text []byte, err error) {
	for _, ch := range instance {
		if !((ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-') {
			return nil, fmt.Errorf("illegal %s: character %c is not supported", instance.TypeName(), ch)
		}
	}
	return []byte(instance), nil
}

func (instance *ResourceIdentifier) UnmarshalText(plain []byte) error {
	return instance.Set(string(plain))
}

func (instance *ResourceIdentifier) UnmarshalJSON(b []byte) error {
	var plain string
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	return instance.Set(plain)
}

func (instance ResourceIdentifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(instance.String())
}

func (instance *ResourceIdentifier) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var plain string
	if err := unmarshal(&plain); err != nil {
		return err
	}
	return instance.Set(plain)
}

func (instance ResourceIdentifier) MarshalYAML() (interface{}, error) {
	return instance.String(), nil
}

func (instance ResourceIdentifier) TypeName() string {
	var result = reflect.TypeOf(instance).Name()
	if len(result) > 0 {
		[]rune(result)[0] = unicode.ToLower([]rune(result)[0])
	}
	return result
}
