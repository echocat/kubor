package model

import (
	"errors"
	"fmt"
)

const (
	DryRunNowhere            = DryRunOn("nowhere")
	DryRunOnClient           = DryRunOn("client")
	DryRunOnServer           = DryRunOn("server")
	DryRunOnServerIfPossible = DryRunOn("serverIfPossible")
)

var (
	ErrIllegalDryRunOn = errors.New("illegal dryRunOn")

	validDryRunOnValues = map[DryRunOn]bool{
		DryRunNowhere:            true,
		DryRunOnClient:           true,
		DryRunOnServer:           true,
		DryRunOnServerIfPossible: true,
	}
)

type DryRunOn string

func (instance *DryRunOn) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance DryRunOn) String() string {
	if exist := validDryRunOnValues[instance]; !exist {
		return fmt.Sprintf("illegal-dry-run-on-%s", string(instance))
	}
	return string(instance)
}

func (instance DryRunOn) MarshalText() (text []byte, err error) {
	if exist := validDryRunOnValues[instance]; !exist {
		return nil, fmt.Errorf("%w: %s", ErrIllegalDryRunOn, string(instance))
	}
	return []byte(instance), nil
}

func (instance *DryRunOn) UnmarshalText(text []byte) error {
	if exist := validDryRunOnValues[DryRunOn(text)]; !exist {
		return fmt.Errorf("%w: %s", ErrIllegalDryRunOn, string(text))
	}
	*instance = DryRunOn(text)
	return nil
}
