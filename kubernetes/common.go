package kubernetes

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strconv"
	"strings"
)

type Project interface {
	GetGroupId() string
	GetArtifactId() string
	GetRelease() string
}

type Object interface {
	runtime.Object
	GroupVersionKind() schema.GroupVersionKind
	GetName() string
	GetNamespace() string
}

type Unstructured interface {
	Object
	UnstructuredContent() map[string]interface{}
}

type Kubor interface {
	GetVersion() string
}

func GroupVersionKindToTypeMeta(kind schema.GroupVersionKind) metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       kind.Kind,
		APIVersion: kind.Version,
	}
}

func Pint32(v int32) *int32 {
	return &v
}

func TryCastToInt32(value interface{}) *int32 {
	switch v := value.(type) {
	case int:
		return Pint32(int32(v))
	case int8:
		return Pint32(int32(v))
	case int16:
		return Pint32(int32(v))
	case int32:
		return Pint32(v)
	case int64:
		return Pint32(int32(v))
	case uint:
		return Pint32(int32(v))
	case uint8:
		return Pint32(int32(v))
	case uint16:
		return Pint32(int32(v))
	case uint32:
		return Pint32(int32(v))
	case uint64:
		return Pint32(int32(v))
	case *int:
		return Pint32(int32(*v))
	case *int8:
		return Pint32(int32(*v))
	case *int16:
		return Pint32(int32(*v))
	case *int32:
		return Pint32(*v)
	case *int64:
		return Pint32(int32(*v))
	case *uint:
		return Pint32(int32(*v))
	case *uint8:
		return Pint32(int32(*v))
	case *uint16:
		return Pint32(int32(*v))
	case *uint32:
		return Pint32(int32(*v))
	case *uint64:
		return Pint32(int32(*v))
	case bool:
		if v {
			return Pint32(1)
		}
		return Pint32(0)
	case *bool:
		if *v {
			return Pint32(1)
		}
		return Pint32(0)
	case string:
		if v, err := strconv.ParseInt(v, 10, 32); err == nil {
			return Pint32(int32(v))
		}
	}
	return Pint32(int32(0))
}

func NormalizeGroupVersionKind(in schema.GroupVersionKind) schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   strings.ToLower(in.Group),
		Version: strings.ToLower(in.Version),
		Kind:    strings.ToLower(in.Kind),
	}
}

func NormalizeGroupVersionKinds(in []schema.GroupVersionKind) []schema.GroupVersionKind {
	result := make([]schema.GroupVersionKind, len(in))
	for i, val := range in {
		result[i] = NormalizeGroupVersionKind(val)
	}
	return result
}

func FormatGroupVersionKind(in schema.GroupVersionKind) string {
	toFormat := NormalizeGroupVersionKind(in)
	result := toFormat.Kind
	if toFormat.Version != "" {
		result = toFormat.Version + "." + result
	}
	if toFormat.Group != "" {
		result = toFormat.Group + "/" + result
	}
	return result
}
