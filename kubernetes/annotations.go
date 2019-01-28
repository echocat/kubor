package kubernetes

import (
	"github.com/levertonai/kubor/kubernetes/fixes"
	"github.com/levertonai/kubor/runtime"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

const (
	AnnotationKuborPrefix            = "kubor.leverton.ai"
	AnnotationKuborVersion           = AnnotationKuborPrefix + "/version"
	AnnotationKuborProjectGroupId    = AnnotationKuborPrefix + "/project-group-id"
	AnnotationKuborProjectArtifactId = AnnotationKuborPrefix + "/project-artifact-id"
	AnnotationKuborProjectRelease    = AnnotationKuborPrefix + "/project-release"
	AnnotationKuborLastAppliedAt     = AnnotationKuborPrefix + "/lastAppliedAt"
)

func init() {
	fixes.ApplyKuborAnnotations = func(object v1.Object, pih fixes.Project) error {
		return SetKuborAnnotations(object, pih)
	}
}

func SetKuborAnnotations(object v1.Object, project Project) error {
	annotations := object.GetAnnotations()

	if annotations == nil {
		annotations = make(map[string]string)
	}

	annotations[AnnotationKuborVersion] = runtime.Runtime.Version
	annotations[AnnotationKuborProjectGroupId] = project.GetGroupId()
	annotations[AnnotationKuborProjectArtifactId] = project.GetArtifactId()
	annotations[AnnotationKuborProjectRelease] = project.GetRelease()
	annotations[AnnotationKuborLastAppliedAt] = time.Now().Format(time.RFC3339)

	object.SetAnnotations(annotations)
	return nil
}

func GetKuborAnnotations(object v1.Object) (KuborAnnotations, error) {
	annotations := object.GetAnnotations()
	if annotations == nil {
		return annotationBasedBasedKuborAnnotations{}, nil
	}
	return annotationBasedBasedKuborAnnotations(annotations), nil
}

type KuborAnnotations interface {
	GetProject() Project
	GetKubor() Kubor
	GetLastAppliedAt() *time.Time
}

type annotationBasedBasedKuborAnnotations map[string]string

func (instance annotationBasedBasedKuborAnnotations) GetProject() Project {
	return annotationBasedBasedProject(instance)
}

func (instance annotationBasedBasedKuborAnnotations) GetKubor() Kubor {
	return annotationBasedBasedKubor(instance)
}

func (instance annotationBasedBasedKuborAnnotations) GetLastAppliedAt() *time.Time {
	plain := instance[AnnotationKuborLastAppliedAt]
	if plain == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, plain)
	if err != nil {
		return nil
	}
	return &t
}

type annotationBasedBasedProject map[string]string

func (instance annotationBasedBasedProject) GetGroupId() string {
	return instance[AnnotationKuborProjectGroupId]
}

func (instance annotationBasedBasedProject) GetArtifactId() string {
	return instance[AnnotationKuborProjectArtifactId]
}

func (instance annotationBasedBasedProject) GetRelease() string {
	return instance[AnnotationKuborProjectRelease]
}

type annotationBasedBasedKubor map[string]string

func (instance annotationBasedBasedKubor) GetVersion() string {
	return instance[AnnotationKuborVersion]
}
