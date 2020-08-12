package common

import (
	"encoding/json"
	"fmt"
	"github.com/echocat/kubor/template"
	"github.com/echocat/kubor/template/functions"
	"reflect"
	"regexp"
	"strings"
)

type EvaluatingPredicate struct {
	Includes EvaluatingPartPredicate
	Excludes EvaluatingPartPredicate
}

func (instance EvaluatingPredicate) IsRelevant() bool {
	return instance.Includes.IsRelevant() || instance.Excludes.IsRelevant()
}

func (instance EvaluatingPredicate) Matches(data interface{}) (bool, error) {
	if len(instance.Includes) > 0 {
		if match, err := instance.Includes.Matches(data); err != nil {
			return false, err
		} else if !match {
			return false, nil
		}
	}
	if len(instance.Excludes) > 0 {
		if match, err := instance.Excludes.Matches(data); err != nil {
			return false, err
		} else if match {
			return false, nil
		}
	}
	return true, nil
}

func (instance *EvaluatingPredicate) UnmarshalJSON(b []byte) error {
	var plains []string
	if err := json.Unmarshal(b, &plains); err != nil {
		return err
	}
	return instance.EvaluatePatterns(plains)
}

func (instance EvaluatingPredicate) MarshalJSON() ([]byte, error) {
	return json.Marshal(instance.Patterns())
}

func (instance *EvaluatingPredicate) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var plains interface{}
	if err := unmarshal(&plains); err != nil {
		return err
	}
	if plains == nil {
		return instance.EvaluatePatterns([]string{})
	}
	if slice, ok := plains.([]string); ok {
		return instance.EvaluatePatterns(slice)
	}
	if str, ok := plains.(string); ok {
		if str == "" {
			return instance.EvaluatePatterns([]string{})
		}
		return instance.EvaluatePatterns([]string{str})
	} else {
		return fmt.Errorf("cannot handle value of type %v for a predicate", reflect.TypeOf(plains))
	}
}

func (instance EvaluatingPredicate) MarshalYAML() (interface{}, error) {
	patterns := instance.Patterns()
	if len(patterns) == 0 {
		return nil, nil
	}
	if len(patterns) == 1 {
		return patterns[0], nil
	}
	return instance.Patterns(), nil
}

func (instance EvaluatingPredicate) String() string {
	var result string
	for i, pattern := range instance.Patterns() {
		if i > 0 {
			result += ", "
		}
		result += pattern
	}
	return result
}

func (instance *EvaluatingPredicate) Set(plain string) error {
	if instance == nil {
		*instance = EvaluatingPredicate{}
	}
	shouldMatch := true
	if strings.HasPrefix(plain, "!") {
		plain = plain[1:]
		shouldMatch = false
	}

	if matcher, err := ParseEvaluatingPartMatcher(plain); err != nil {
		return err
	} else if shouldMatch {
		instance.Includes = append(instance.Includes, matcher)
	} else {
		instance.Excludes = append(instance.Excludes, matcher)
	}
	return nil
}

func (instance EvaluatingPredicate) Patterns() []string {
	patterns := make([]string, len(instance.Includes)+len(instance.Excludes))
	for i, part := range instance.Includes {
		patterns[i] = part.String()
	}
	for i, part := range instance.Excludes {
		patterns[i+len(instance.Includes)] = fmt.Sprintf("!%s", part.String())
	}
	return patterns
}

func (instance *EvaluatingPredicate) EvaluatePatterns(patterns []string) error {
	result := EvaluatingPredicate{}
	for _, pattern := range patterns {
		if err := result.Set(pattern); err != nil {
			return err
		}
	}
	*instance = result
	return nil
}

func ParseEvaluatingPartMatcher(plain string) (EvaluatingPartMatcher, error) {
	parts := strings.SplitN(plain, "=", 2)
	if len(parts) != 2 {
		return EvaluatingPartMatcher{}, fmt.Errorf("illegal matching pattern '%s': '=' missing", plain)
	} else if matcher, err := NewEvaluatingPartMatcher(parts[0], parts[1]); err != nil {
		return EvaluatingPartMatcher{}, fmt.Errorf("illegal matching pattern '%s': %w", plain, err)
	} else {
		return matcher, nil
	}
}

func NewEvaluatingPartMatcher(valueTemplate string, check string) (EvaluatingPartMatcher, error) {
	if tmpl, err := functions.DefaultTemplateFactory().New("", valueTemplate); err != nil {
		return EvaluatingPartMatcher{}, err
	} else if checkInstance, err := regexp.Compile(fmt.Sprintf("^%s$", check)); err != nil {
		return EvaluatingPartMatcher{}, err
	} else {
		return EvaluatingPartMatcher{
			valueTemplateSource: valueTemplate,
			valueTemplate:       tmpl,
			check:               checkInstance,
		}, nil
	}
}

type EvaluatingPartMatcher struct {
	valueTemplateSource string
	valueTemplate       template.Template
	check               *regexp.Regexp
}

func (instance EvaluatingPartMatcher) IsRelevant() bool {
	return true
}

func (instance EvaluatingPartMatcher) Value(data interface{}) (string, error) {
	return instance.valueTemplate.ExecuteToString(data)
}

func (instance EvaluatingPartMatcher) Matches(data interface{}) (bool, error) {
	if val, err := instance.Value(data); err != nil {
		return false, err
	} else {
		return instance.check.MatchString(val), nil
	}
}

func (instance EvaluatingPartMatcher) String() string {
	return fmt.Sprintf("%s=%v", instance.valueTemplateSource, instance.check)
}

type EvaluatingPartPredicate []EvaluatingPartMatcher

func (instance EvaluatingPartPredicate) IsRelevant() bool {
	for _, candidate := range instance {
		if candidate.IsRelevant() {
			return true
		}
	}
	return false
}

func (instance EvaluatingPartPredicate) Matches(data interface{}) (bool, error) {
	for _, matcher := range instance {
		if match, err := matcher.Matches(data); err != nil {
			return false, err
		} else if match {
			return true, nil
		}
	}
	return false, nil
}
