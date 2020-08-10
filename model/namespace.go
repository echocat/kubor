package model

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	namespaceRegexp = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

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
	if len(instance) > 0 && (!namespaceRegexp.MatchString(string(instance)) || len(instance) > 253) {
		return []byte(fmt.Sprintf("illegal-namespace-%s", string(instance))),
			fmt.Errorf("%w: %s", ErrIllegalNamespace, string(instance))
	}
	return []byte(instance), nil
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
