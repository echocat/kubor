package model

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	nameRegexp = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

	ErrIllegalName = errors.New("illegal name")
)

type Name string

func (instance *Name) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance Name) String() string {
	v, _ := instance.MarshalText()
	return string(v)
}

func (instance Name) MarshalText() (text []byte, err error) {
	if len(instance) > 0 && (!nameRegexp.MatchString(string(instance)) || len(instance) > 253) {
		return []byte(fmt.Sprintf("illega-name-%s", string(instance))),
			fmt.Errorf("%w: %s", ErrIllegalName, string(instance))
	}
	return []byte(instance), nil
}

func (instance *Name) UnmarshalText(text []byte) error {
	v := Name(text)
	if _, err := instance.MarshalText(); err != nil {
		return err
	}
	*instance = v
	return nil
}
