package model

import (
	"errors"
	"fmt"
	"strings"
)

type LabelAction uint8

var (
	ErrIllegalLabelAction = errors.New("illegal label-action")
)

const (
	LabelActionSetIfAbsent = LabelAction(0)
	LabelActionSet         = LabelAction(1)
	LabelActionSetIfExists = LabelAction(2)
	LabelActionLeave       = LabelAction(3)
	LabelActionDrop        = LabelAction(4)
)

func (instance *LabelAction) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance LabelAction) String() string {
	v, _ := instance.MarshalText()
	return string(v)
}

func (instance LabelAction) MarshalText() (text []byte, err error) {
	switch instance {
	case LabelActionSetIfAbsent:
		return []byte("set-if-absent"), nil
	case LabelActionSet:
		return []byte("set"), nil
	case LabelActionSetIfExists:
		return []byte("set-if-exists"), nil
	case LabelActionLeave:
		return []byte("leave"), nil
	case LabelActionDrop:
		return []byte("drop"), nil
	default:
		return []byte(fmt.Sprintf("illegal-label-action-%d", instance)),
			fmt.Errorf("%w: %s", ErrIllegalLabelAction, string(instance))
	}
}

func (instance *LabelAction) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "set-if-absent", "setifabsent", "":
		*instance = LabelActionSetIfAbsent
		return nil
	case "set":
		*instance = LabelActionSet
		return nil
	case "set-if-exists", "set-if-exist", "setifexists", "setifexist":
		*instance = LabelActionSetIfExists
		return nil
	case "leave", "ignore":
		*instance = LabelActionLeave
		return nil
	case "drop":
		*instance = LabelActionDrop
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrIllegalLabelAction, string(text))
	}
}
