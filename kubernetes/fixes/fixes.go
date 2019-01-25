package fixes

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	updateFixes []updateFix
	createFixes []createFix
)

type Project interface {
	GetGroupId() string
	GetArtifactId() string
	GetRelease() string
}

type updateFix func(project Project, original unstructured.Unstructured, target *unstructured.Unstructured) error

type createFix func(project Project, target *unstructured.Unstructured) error

func registerUpdateFix(fix updateFix) {
	updateFixes = append(updateFixes, fix)
}

func registerCreateFix(fix createFix) {
	createFixes = append(createFixes, fix)
}

func FixForUpdate(project Project, original unstructured.Unstructured, target *unstructured.Unstructured) error {
	for _, fix := range updateFixes {
		if err := fix(project, original, target); err != nil {
			return err
		}
	}
	return nil
}

func FixForCreate(project Project, target *unstructured.Unstructured) error {
	for _, fix := range createFixes {
		if err := fix(project, target); err != nil {
			return err
		}
	}
	return nil
}
