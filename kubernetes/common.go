package kubernetes

import (
	"github.com/echocat/kubor/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"strings"
)

func GroupVersionKindToTypeMeta(kind model.GroupVersionKind) metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       strings.ToLower(kind.Kind),
		APIVersion: strings.ToLower(kind.Version),
	}
}

func Pbool(v bool) *bool {
	return &v
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
