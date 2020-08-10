package model

import (
	"errors"
	"fmt"
	"strings"
)

type DryRun uint8

const (
	DryRunBefore = DryRun(0)
	DryRunOnly   = DryRun(1)
	DryRunNever  = DryRun(2)
)

var (
	ErrIllegalDryRun = errors.New("illegal dry-Run")
)

func (instance *DryRun) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance DryRun) String() string {
	v, _ := instance.MarshalText()
	return string(v)
}

func (instance DryRun) MarshalText() (text []byte, err error) {
	switch instance {
	case DryRunBefore:
		return []byte("before"), nil
	case DryRunOnly:
		return []byte("only"), nil
	case DryRunNever:
		return []byte("never"), nil
	default:
		return []byte(fmt.Sprintf("illegal-dry-Run-%d", instance)),
			fmt.Errorf("%w: %d", ErrIllegalDryRun, instance)
	}
}

func (instance *DryRun) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "before", "", "default":
		*instance = DryRunBefore
		return nil
	case "only":
		*instance = DryRunOnly
		return nil
	case "never", "off", "false":
		*instance = DryRunNever
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrIllegalDryRun, string(text))
	}
}

func (instance DryRun) IsDryRunAllowed() bool {
	switch instance {
	case DryRunBefore, DryRunOnly:
		return true
	default:
		return false
	}
}

func (instance DryRun) IsApplyAllowed() bool {
	switch instance {
	case DryRunBefore, DryRunNever:
		return true
	default:
		return false
	}
}
