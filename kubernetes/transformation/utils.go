package transformation

import (
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
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

func NestedNamedSliceAsMaps(obj map[string]interface{}, nameField string, fields ...string) (result map[string]map[string]interface{}, found bool, err error) {
	slice, found, err := unstructured.NestedSlice(obj, fields...)
	if err != nil || !found {
		return map[string]map[string]interface{}{}, found, err
	}
	result, err = sliceToNamedMap(slice, nameField)
	return
}

func SetNestedNamedMapsAsSlice(obj map[string]interface{}, nameField string, v map[string]map[string]interface{}, fields ...string) error {
	slice := mapToNamedSlice(v, nameField)
	return unstructured.SetNestedSlice(obj, slice, fields...)
}

func sliceToNamedMap(in []interface{}, nameField string) (result map[string]map[string]interface{}, err error) {
	result = make(map[string]map[string]interface{}, len(in))

	for _, entry := range in {
		if m, ok := entry.(map[string]interface{}); !ok {
			return nil, fmt.Errorf("expected entry of type map[string]interface{}, but got: %v", reflect.TypeOf(entry))
		} else if vName, ok := m[nameField]; !ok {
			result[""] = m
		} else {
			delete(m, nameField)
			result[fmt.Sprint(vName)] = m
		}
	}

	return
}

func mapToNamedSlice(in map[string]map[string]interface{}, nameField string) (result []interface{}) {
	result = make([]interface{}, len(in))

	var i int
	for name, entry := range in {
		entry[nameField] = name
		result[i] = entry
		i++
	}

	return
}
