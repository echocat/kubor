package model

import (
	"errors"
	"fmt"
	"strings"
)

const (
	ApplyOnAlways = ApplyOn("always")
	ApplyOnCreate = ApplyOn("create")
	ApplyOnUpdate = ApplyOn("update")
	ApplyOnNever  = ApplyOn("never")
)

var (
	ErrIllegalApplyOn = errors.New("illegal apply-on")

	validApplyOnValues = map[ApplyOn]bool{ApplyOnAlways: true, ApplyOnCreate: true, ApplyOnUpdate: true, ApplyOnNever: true}
)

type ApplyOn string

func (instance *ApplyOn) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance ApplyOn) String() string {
	if exist := validApplyOnValues[instance]; !exist {
		return fmt.Sprintf("illegal-apply-on-%s", string(instance))
	}
	return string(instance)
}

func (instance ApplyOn) MarshalText() (text []byte, err error) {
	if exist := validApplyOnValues[instance]; !exist {
		return nil, fmt.Errorf("%w: %s", ErrIllegalApplyOn, string(instance))
	}
	return []byte(instance), nil
}

func (instance ApplyOn) OnUpdate() bool {
	switch instance {
	case ApplyOnUpdate, ApplyOnAlways:
		return true
	default:
		return false
	}
}

func (instance ApplyOn) OnCreate() bool {
	switch instance {
	case ApplyOnCreate, ApplyOnAlways:
		return true
	default:
		return false
	}
}

func (instance *ApplyOn) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "true", "enabled", "on":
		*instance = ApplyOnAlways
		return nil
	case "false", "disabled", "off":
		*instance = ApplyOnNever
		return nil
	}

	if exist := validApplyOnValues[ApplyOn(text)]; !exist {
		return fmt.Errorf("%w: %s", ErrIllegalApplyOn, string(text))
	}
	*instance = ApplyOn(text)
	return nil
}
