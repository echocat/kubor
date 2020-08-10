package transformation

import (
	"k8s.io/apimachinery/pkg/runtime"
	"strings"
)

func groupVersionKindMatchesVersion(what runtime.Object, candidates ...string) bool {
	if what == nil {
		return false
	}
	kind := what.GetObjectKind()
	if kind == nil {
		return false
	}
	target := strings.ToLower(kind.GroupVersionKind().Version)
	for _, candidate := range candidates {
		if strings.ToLower(candidate) == target {
			return true
		}
	}
	return false
}

func groupVersionKindMatchesKind(what runtime.Object, candidates ...string) bool {
	if what == nil {
		return false
	}
	kind := what.GetObjectKind()
	if kind == nil {
		return false
	}
	target := strings.ToLower(kind.GroupVersionKind().Kind)
	for _, candidate := range candidates {
		if strings.ToLower(candidate) == target {
			return true
		}
	}
	return false
}

func groupVersionKindMatches(left, right runtime.Object) bool {
	if left == nil || right == nil {
		return false
	}
	kindLeft := left.GetObjectKind()
	kindRight := right.GetObjectKind()
	if kindLeft == nil || kindRight == nil {
		return false
	}
	gvkLeft := kindLeft.GroupVersionKind()
	gvkRight := kindRight.GroupVersionKind()
	if strings.ToLower(gvkLeft.Group) != strings.ToLower(gvkRight.Group) {
		return false
	}
	if strings.ToLower(gvkLeft.Version) != strings.ToLower(gvkRight.Version) {
		return false
	}
	if strings.ToLower(gvkLeft.Kind) != strings.ToLower(gvkRight.Kind) {
		return false
	}
	return true
}
