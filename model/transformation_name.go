package model

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	transformationNameRegexp = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

	ErrIllegalTransformationName = errors.New("illegal transformation name")
)

type TransformationName Name

func (instance *TransformationName) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance TransformationName) String() string {
	v, _ := instance.MarshalText()
	return string(v)
}

func (instance TransformationName) MarshalText() (text []byte, err error) {
	if len(instance) > 0 && (!transformationNameRegexp.MatchString(string(instance)) || len(instance) > 253) {
		return []byte(fmt.Sprintf("illegal-transformation-name-%s", string(instance))),
			fmt.Errorf("%w: %s", ErrIllegalTransformationName, string(instance))
	}
	return []byte(instance), nil
}

func (instance *TransformationName) UnmarshalText(text []byte) error {
	v := TransformationName(text)
	if _, err := instance.MarshalText(); err != nil {
		return err
	}
	*instance = v
	return nil
}
