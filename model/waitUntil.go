package model

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"time"
)

var (
	ErrIllegalWaitUntil      = errors.New("illegal wait-until")
	ErrIllegalWaitUntilStage = errors.New("illegal wait-until-stage")

	WaitUntilDefault = WaitUntil{
		Stage: WaitUntilStageDefault,
	}
)

type WaitUntil struct {
	Stage   WaitUntilStage
	Timeout *time.Duration
}

func (instance WaitUntil) CopyWithTimeout(timeout *time.Duration) WaitUntil {
	return WaitUntil{
		Stage:   instance.Stage,
		Timeout: timeout,
	}
}

func (instance WaitUntil) ShouldInherit() bool {
	return instance.Stage == WaitUntilStageDefault
}

func (instance *WaitUntil) Set(plain string) error {
	if duration, err := time.ParseDuration(plain); err == nil {
		*instance = WaitUntil{WaitUntilStageApplied, &duration}
		return nil
	}

	parts := strings.Split(plain, ":")
	var stage WaitUntilStage
	if err := stage.Set(parts[0]); err != nil {
		return fmt.Errorf("%w, %v: %s", ErrIllegalWaitUntil, err, plain)
	}
	result := WaitUntil{Stage: stage}

	if len(parts) > 2 {
		return fmt.Errorf("%w, expected no or one duration argument but got %d: %s", ErrIllegalWaitUntil, len(parts)-1, plain)
	}
	if len(parts) == 2 {
		if d, err := time.ParseDuration(parts[1]); err != nil {
			return fmt.Errorf("%w, %v: %s", ErrIllegalWaitUntil, err, plain)
		} else if d > 0 {
			result.Timeout = &d
		}
	}
	*instance = result
	return nil
}

func (instance WaitUntil) String() string {
	if t := instance.Timeout; t != nil {
		return instance.Stage.String() + ":" + t.String()
	}
	return instance.Stage.String()
}

func (instance WaitUntil) MarshalText() (text []byte, err error) {
	return []byte(instance.String()), nil
}

func (instance *WaitUntil) UnmarshalText(text []byte) error {
	return instance.Set(string(text))
}

func (instance WaitUntil) AsLazyFormatter(template string) *WaitUntilLazyFormatter {
	return &WaitUntilLazyFormatter{
		WaitUntil: instance,
		Template:  template,
	}
}

func (instance WaitUntil) MergeWith(child WaitUntil) (result WaitUntil) {
	if child.Stage == WaitUntilStageDefault {
		return instance
	}
	return child
}

type WaitUntilStage uint8

const (
	WaitUntilStageDefault  = WaitUntilStage(0)
	WaitUntilStageNever    = WaitUntilStage(1)
	WaitUntilStageApplied  = WaitUntilStage(2)
	WaitUntilStageDeployed = WaitUntilStage(3)
	WaitUntilStageExecuted = WaitUntilStage(4)
)

func (instance WaitUntilStage) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance WaitUntilStage) String() string {
	str, _ := instance.MarshalText()
	return string(str)
}

func (instance WaitUntilStage) MarshalText() (text []byte, err error) {
	switch instance {
	case WaitUntilStageDefault:
		return []byte("default"), nil
	case WaitUntilStageNever:
		return []byte("never"), nil
	case WaitUntilStageApplied:
		return []byte("applied"), nil
	case WaitUntilStageDeployed:
		return []byte("deployed"), nil
	case WaitUntilStageExecuted:
		return []byte("executed"), nil
	default:
		return []byte(fmt.Sprintf("illegal-wait-until-stage-%d", instance)),
			fmt.Errorf("%w: %d", ErrIllegalWaitUntilStage, instance)
	}
}

func (instance *WaitUntilStage) UnmarshalText(text []byte) error {
	switch string(text) {
	case "default", "inherit", "":
		*instance = WaitUntilStageDefault
		return nil
	case "never":
		*instance = WaitUntilStageNever
		return nil
	case "applied", "apply":
		*instance = WaitUntilStageApplied
		return nil
	case "deployed", "deploy":
		*instance = WaitUntilStageDeployed
		return nil
	case "executed", "execute":
		*instance = WaitUntilStageExecuted
		return nil
	default:
		return fmt.Errorf("%w: %d", ErrIllegalWaitUntilStage, instance)
	}
}

type WaitUntilLazyFormatter struct {
	WaitUntil
	Template string
}

func (instance WaitUntilLazyFormatter) String() string {
	text, err := instance.MarshalText()
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(text)
}

func (instance WaitUntilLazyFormatter) MarshalText() (text []byte, err error) {
	tmpl, err := template.New(instance.Template).Parse(instance.Template)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, instance.WaitUntil); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
