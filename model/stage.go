package model

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	stageRegexp = regexp.MustCompile(`^[a-z][a-z0-9.:_-]*[a-z0-9]$`)

	ErrIllegalStage      = errors.New("illegal stage")
	ErrIllegalStageRange = errors.New("illegal stage-range")
)

const (
	StageDefault = Stage("deploy")
)

type Stage string

func (instance *Stage) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance Stage) String() string {
	if !stageRegexp.MatchString(string(instance)) {
		return fmt.Sprintf("illegal-stage-%s", string(instance))
	}
	return string(instance)
}

func (instance Stage) MarshalText() (text []byte, err error) {
	if !stageRegexp.MatchString(string(instance)) {
		return []byte(fmt.Sprintf("illega-stage-%s", string(instance))),
			fmt.Errorf("%w: %s", ErrIllegalStage, string(instance))
	}
	return []byte(instance), nil
}

func (instance *Stage) UnmarshalText(text []byte) error {
	if !stageRegexp.MatchString(string(text)) {
		return fmt.Errorf("%w: %s", ErrIllegalStage, string(text))
	}
	*instance = Stage(text)
	return nil
}

type Stages []Stage

func (instance Stages) Contains(what Stage) bool {
	if len(instance) == 0 && what == StageDefault {
		return true
	}
	for _, candidate := range instance {
		if candidate == what {
			return true
		}
	}
	return false
}

func (instance *Stages) Set(plain string) error {
	result := Stages{}
	for _, plainPart := range strings.Split(plain, ",") {
		plainPart = strings.TrimSpace(plainPart)
		var part Stage
		if err := part.Set(plainPart); err != nil {
			return err
		}
		result = append(result, part)
	}
	*instance = result
	return nil
}

func (instance Stages) Strings() []string {
	result := make([]string, len(instance))
	for i, part := range instance {
		result[i] = part.String()
	}
	return result
}

func (instance Stages) String() string {
	result := strings.Join(instance.Strings(), ",")
	if result == "" {
		return StageDefault.String()
	}
	return result
}

type StageRange struct {
	From *Stage
	To   *Stage
}

func (instance StageRange) IsRelevant() bool {
	return instance.From != nil || instance.To != nil
}

func (instance StageRange) Matches(stages Stages, stage Stage) bool {
	if len(stages) == 0 && stage == StageDefault {
		// Default behavior if no stages are defined and the stage is the default one.
		return true
	}

	started := false
	for _, current := range stages {
		if !started {
			if instance.From == nil || current == *instance.From {
				started = true
			} else {
				continue
			}
		}

		if current == stage {
			return true
		}

		if instance.To != nil && current == *instance.To {
			return false
		}
	}
	return false
}

func (instance *StageRange) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance StageRange) String() string {
	v, _ := instance.MarshalText()
	return string(v)
}

func (instance StageRange) MarshalText() (text []byte, err error) {
	from, to := "", ""
	if s := instance.From; s != nil {
		if v, err := s.MarshalText(); err != nil {
			return v, err
		} else {
			from = string(v)
		}
	}
	if s := instance.To; s != nil {
		if v, err := s.MarshalText(); err != nil {
			return v, err
		} else {
			from = string(v)
		}
	}
	return []byte(from + ":" + to), nil
}

func (instance *StageRange) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*instance = StageRange{}
		return nil
	}
	parts := strings.Split(string(text), ":")
	if len(parts) > 2 {
		return fmt.Errorf("%w: expected 1 or 2 parts separated by a ':', but got: %d", ErrIllegalStageRange, len(parts))
	}
	result := StageRange{}

	if parts[0] != "" {
		var v Stage
		if err := v.UnmarshalText([]byte(parts[0])); err != nil {
			return fmt.Errorf("%w: %v", ErrIllegalStageRange, err)
		}
		result.From = &v
	}

	if len(parts) == 1 {
		result.To = result.From
	} else if parts[1] != "" {
		var v Stage
		if err := v.UnmarshalText([]byte(parts[1])); err != nil {
			return fmt.Errorf("%w: %v", ErrIllegalStageRange, err)
		}
		result.To = &v
	}

	*instance = result
	return nil
}
