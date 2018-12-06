package fixes

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	updateFixes []updateFix
	createFixes []createFix
)

type updateFix func(original unstructured.Unstructured, target *unstructured.Unstructured) error

type createFix func(target *unstructured.Unstructured) error

func registerUpdateFix(fix updateFix) {
	updateFixes = append(updateFixes, fix)
}

func registerCreateFix(fix createFix) {
	createFixes = append(createFixes, fix)
}

func FixForUpdate(original unstructured.Unstructured, target *unstructured.Unstructured) error {
	for _, fix := range updateFixes {
		if err := fix(original, target); err != nil {
			return err
		}
	}
	return nil
}

func FixForCreate(target *unstructured.Unstructured) error {
	for _, fix := range createFixes {
		if err := fix(target); err != nil {
			return err
		}
	}
	return nil
}
