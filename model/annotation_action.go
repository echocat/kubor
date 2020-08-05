package model

import (
	"errors"
	"fmt"
	"strings"
)

type AnnotationAction uint8

var (
	ErrIllegalAnnotationAction = errors.New("illegal annotation-action")
)

const (
	AnnotationActionDrop  = AnnotationAction(0)
	AnnotationActionLeave = AnnotationAction(1)
)

func (instance *AnnotationAction) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance AnnotationAction) String() string {
	v, _ := instance.MarshalText()
	return string(v)
}

func (instance AnnotationAction) MarshalText() (text []byte, err error) {
	switch instance {
	case AnnotationActionDrop:
		return []byte("drop"), nil
	case AnnotationActionLeave:
		return []byte("leave"), nil
	default:
		return []byte(fmt.Sprintf("illegal-annotation-action-%d", instance)),
			fmt.Errorf("%w: %s", ErrIllegalAnnotationAction, string(instance))
	}
}

func (instance *AnnotationAction) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "drop", "":
		*instance = AnnotationActionDrop
		return nil
	case "leave", "ignore":
		*instance = AnnotationActionLeave
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrIllegalAnnotationAction, string(text))
	}
}
