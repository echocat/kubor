package model

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

const (
	AnnotationStage     = "kubor.echocat.org/stage"
	AnnotationApplyOn   = "kubor.echocat.org/apply-on"
	AnnotationWaitUntil = "kubor.echocat.org/wait-until"
)

type Annotations struct {
	Stage     Annotation `yaml:"stage,omitempty" json:"stage,omitempty"`
	ApplyOn   Annotation `yaml:"applyOn,omitempty" json:"applyOn,omitempty"`
	WaitUntil Annotation `yaml:"waitUntil,omitempty" json:"waitUntil,omitempty"`
}

func newAnnotations() Annotations {
	return Annotations{
		Stage:     Annotation{AnnotationStage, AnnotationActionDrop},
		ApplyOn:   Annotation{AnnotationApplyOn, AnnotationActionDrop},
		WaitUntil: Annotation{AnnotationWaitUntil, AnnotationActionDrop},
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

func (instance Annotations) GetWaitUntilFor(v *unstructured.Unstructured) (WaitUntil, error) {
	as := v.GetAnnotations()
	plain := as[string(instance.ApplyOn.Name)]
	if plain == "" {
		return WaitUntilDefault, nil
	}
	var result WaitUntil
	return result, result.Set(plain)
}
