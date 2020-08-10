package model

import (
	"errors"
	"fmt"
	"regexp"
)

type Label struct {
	Name   MetaName    `yaml:"name" json:"name"`
	Action LabelAction `yaml:"action" json:"action"`
}

var (
	labelValueRegexp = regexp.MustCompile(`^[a-z0-9A-Z]([a-z0-9A-Z._-]*[a-z0-9A-Z]|)$`)

	ErrIllegalLabelValue = errors.New("illegal label-value")
)

type LabelValue string

func (instance *LabelValue) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance LabelValue) String() string {
	if v, err := instance.MarshalText(); err != nil {
		return fmt.Sprintf("illegal-label-value-%s", string(instance))
	} else {
		return string(v)
	}
}

func (instance LabelValue) MarshalText() (text []byte, err error) {
	value := string(instance)
	if value != "" && (!labelValueRegexp.MatchString(value) || len(value) > 63) {
		return nil, fmt.Errorf("%w: %s", ErrIllegalLabelValue, string(instance))
	}

	return []byte(instance), nil
}

func (instance *LabelValue) UnmarshalText(text []byte) error {
	result := LabelValue(text)
	if _, err := result.MarshalText(); err != nil {
		return err
	}
	*instance = result
	return nil
}
