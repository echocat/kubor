package model

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrIllegalNamespace = errors.New("illegal namespace")
)

type Namespace Name

func (instance *Namespace) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance Namespace) String() string {
	v, _ := instance.MarshalText()
	return string(v)
}

func (instance Namespace) MarshalText() (text []byte, err error) {
	if v, err := Name(instance).MarshalText(); err != nil {
		return []byte(fmt.Sprintf("illega-namespace-%s", string(instance))),
			fmt.Errorf("%w: %s", ErrIllegalNamespace, string(instance))
	} else {
		return v, nil
	}
}

func (instance *Namespace) UnmarshalText(text []byte) error {
	v := Namespace(text)
	if _, err := instance.MarshalText(); err != nil {
		return err
	}
	*instance = v
	return nil
}

type Namespaces []Namespace

func (instance Namespaces) Contains(what Namespace) bool {
	if len(instance) == 0 || what == "" {
		return true
	}
	for _, candidate := range instance {
		if candidate == what {
			return true
		}
	}
	return false
}

func (instance *Namespaces) Set(plain string) error {
	result := Namespaces{}
	for _, plainPart := range strings.Split(plain, ",") {
		plainPart = strings.TrimSpace(plainPart)
		var part Namespace
		if err := part.Set(plainPart); err != nil {
			return err
		}
		result = append(result, part)
	}
	*instance = result
	return nil
}

func (instance Namespaces) Strings() []string {
	result := make([]string, len(instance))
	for i, part := range instance {
		result[i] = part.String()
	}
	return result
}

func (instance Namespaces) String() string {
	return strings.Join(instance.Strings(), ",")
}
