package model

type Annotation struct {
	Name   MetaName         `yaml:"name" json:"name"`
	Action AnnotationAction `yaml:"action" json:"action"`
}
