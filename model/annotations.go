package model

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strings"
)

const (
	AnnotationStage                = "kubor.echocat.org/stage"
	AnnotationApplyOn              = "kubor.echocat.org/apply-on"
	AnnotationDryRunOn             = "kubor.echocat.org/dry-run-on"
	AnnotationWaitUntil            = "kubor.echocat.org/wait-until"
	AnnotationCleanupOn            = "kubor.echocat.org/cleanup-on"
	AnnotationTransformationPrefix = "transformation.kubor.echocat.org/"
)

type Annotations struct {
	Stage           Annotation `yaml:"stage,omitempty" json:"stage,omitempty"`
	ApplyOn         Annotation `yaml:"applyOn,omitempty" json:"applyOn,omitempty"`
	DryRunOn        Annotation `yaml:"dryRunOn,omitempty" json:"dryRunOn,omitempty"`
	WaitUntil       Annotation `yaml:"waitUntil,omitempty" json:"waitUntil,omitempty"`
	CleanupOn       Annotation `yaml:"cleanupOn,omitempty" json:"cleanupOn,omitempty"`
	Transformations Annotation `yaml:"transformations,omitempty" json:"transformations,omitempty"`
}

func NewAnnotations() Annotations {
	return Annotations{
		Stage:           Annotation{AnnotationStage, AnnotationActionDrop},
		ApplyOn:         Annotation{AnnotationApplyOn, AnnotationActionDrop},
		DryRunOn:        Annotation{AnnotationDryRunOn, AnnotationActionDrop},
		WaitUntil:       Annotation{AnnotationWaitUntil, AnnotationActionDrop},
		CleanupOn:       Annotation{AnnotationCleanupOn, AnnotationActionLeave},
		Transformations: Annotation{AnnotationTransformationPrefix, AnnotationActionDrop},
	}
}

func (instance Annotations) GetStageFor(v *unstructured.Unstructured) (Stage, error) {
	as := v.GetAnnotations()
	plain := as[string(instance.Stage.Name)]
	if plain == "" {
		return StageDefault, nil
	}
	var result Stage
	return result, result.Set(plain)
}

func (instance Annotations) GetApplyOnFor(v *unstructured.Unstructured) (ApplyOn, error) {
	as := v.GetAnnotations()
	plain := as[string(instance.ApplyOn.Name)]
	if plain == "" {
		return ApplyOnAlways, nil
	}
	var result ApplyOn
	return result, result.Set(plain)
}

func (instance Annotations) GetDryRunOnFor(v *unstructured.Unstructured, def DryRunOn) (DryRunOn, error) {
	as := v.GetAnnotations()
	plain := as[string(instance.DryRunOn.Name)]
	if plain == "" {
		return def, nil
	}
	var result DryRunOn
	return result, result.Set(plain)
}

func (instance Annotations) GetWaitUntilFor(v *unstructured.Unstructured) (WaitUntil, error) {
	as := v.GetAnnotations()
	plain := as[string(instance.WaitUntil.Name)]
	if plain == "" {
		return WaitUntilDefault, nil
	}
	var result WaitUntil
	return result, result.Set(plain)
}

func (instance Annotations) GetCleanupOn(v *unstructured.Unstructured) (CleanupOn, error) {
	as := v.GetAnnotations()
	plain := as[string(instance.CleanupOn.Name)]
	if plain == "" {
		return CleanupOnAutomatic, nil
	}
	var result CleanupOn
	return result, result.Set(plain)
}

func (instance Annotations) GetTransformation(v *unstructured.Unstructured, name TransformationName) (result Transformation, err error) {
	as := v.GetAnnotations()
	plain := as[string(instance.Transformations.Name)+string(name)]
	plain = strings.TrimSpace(plain)
	parts := strings.SplitN(plain, ":", 2)
	switch strings.ToLower(parts[0]) {
	case "":
		// ignore
	case "enabled", "true", "on":
		v := true
		result.Enabled = &v
		if len(parts) > 1 {
			arg := parts[1]
			result.Argument = &arg
		}
	case "disabled", "false", "off":
		v := false
		result.Enabled = &v
		if len(parts) > 1 {
			arg := parts[1]
			result.Argument = &arg
		}
	default:
		result.Argument = &plain
	}
	return
}
