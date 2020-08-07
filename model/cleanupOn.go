package model

import (
	"errors"
	"fmt"
	"strings"
)

type CleanupOn uint8

const (
	CleanupOnAutomatic = CleanupOn(0)
	CleanupOnNever     = CleanupOn(1)
	CleanupOnOrphaned  = CleanupOn(2)
	CleanupOnExecuted  = CleanupOn(3)
	CleanupOnDelete    = CleanupOn(4)
)

var (
	ErrIllegalCleanupOn = errors.New("illegal cleanup-on")
)

func (instance *CleanupOn) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance CleanupOn) String() string {
	v, _ := instance.MarshalText()
	return string(v)
}

func (instance CleanupOn) MarshalText() (text []byte, err error) {
	switch instance {
	case CleanupOnAutomatic:
		return []byte("automatic"), nil
	case CleanupOnNever:
		return []byte("never"), nil
	case CleanupOnOrphaned:
		return []byte("orphaned"), nil
	case CleanupOnExecuted:
		return []byte("executed"), nil
	case CleanupOnDelete:
		return []byte("delete"), nil
	default:
		return []byte(fmt.Sprintf("illegal-cleanup-on-%d", instance)),
			fmt.Errorf("%w: %d", ErrIllegalCleanupOn, instance)
	}
}

func (instance *CleanupOn) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "automatic", "auto", "", "default":
		*instance = CleanupOnAutomatic
		return nil
	case "never", "off", "false":
		*instance = CleanupOnNever
		return nil
	case "orphaned":
		*instance = CleanupOnOrphaned
		return nil
	case "executed":
		*instance = CleanupOnExecuted
		return nil
	case "delete", "deleted", "remove", "removed":
		*instance = CleanupOnDelete
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrIllegalCleanupOn, string(text))
	}
}

func (instance CleanupOn) OnOrphaned() bool {
	switch instance {
	case CleanupOnOrphaned, CleanupOnAutomatic, CleanupOnExecuted:
		return true
	default:
		return false
	}
}

func (instance CleanupOn) OnExecuted() bool {
	switch instance {
	case CleanupOnExecuted, CleanupOnAutomatic:
		return true
	default:
		return false
	}
}

func (instance CleanupOn) OnPurge() bool {
	switch instance {
	case CleanupOnNever:
		return false
	default:
		return true
	}
}
