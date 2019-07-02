package model

import (
	"bytes"
	"fmt"
	"github.com/echocat/kubor/common"
	"io"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/runtime/serializer/streaming"
	"k8s.io/client-go/kubernetes/scheme"
	"strings"
)

// ContentProvider provides the next resource and returns an error of io.EOF if no more element is available.
type ContentProvider func() (name string, content []byte, err error)

type OnObject func(source string, object runtime.Object, unstructured *unstructured.Unstructured) error

func NewObjectHandler(onObject OnObject) (*ObjectHandler, error) {
	return &ObjectHandler{
		OnObject:     onObject,
		Deserializer: scheme.Codecs.UniversalDeserializer(),
	}, nil
}

type ObjectHandler struct {
	OnObject OnObject

	Deserializer runtime.Decoder
}

func (instance *ObjectHandler) Handle(cp ContentProvider) error {
	var name string
	var content []byte
	var err error
	for name, content, err = cp(); err == nil; name, content, err = cp() {
		if err := instance.handleContent(name, content); err != nil {
			return err
		}
	}
	if err == io.EOF {
		return nil
	}
	if se, ok := err.(*errors.StatusError); ok {
		return se
	}
	return fmt.Errorf("cannot handle '%s': %v", name, err)
}

func (instance *ObjectHandler) handleContent(source string, content []byte) error {
	plain := strings.TrimSpace(string(content))
	plain = strings.Replace(plain, "\r\n", "\n", -1)
	if strings.HasPrefix(plain, "---\n") {
		plain = plain[5:]
	}
	if strings.HasSuffix(plain, "\n---") {
		plain = plain[:5]
	}
	parts := strings.Split(plain, "\n---\n")
	for i, part := range parts {
		if strings.TrimSpace(part) != "" {
			fSource := fmt.Sprintf("%s#%d", source, i)
			if object, _, err := instance.Deserializer.Decode([]byte(part), nil, nil); err != nil {
				return fmt.Errorf("%s: %v", fSource, err)
			} else if unstr, err := instance.decodeUnstructured([]byte(part)); err != nil {
				return fmt.Errorf("%s: %v", fSource, err)
			} else if err := instance.OnObject(fSource, object, unstr); err != nil {
				return fmt.Errorf("%s: %v", fSource, err)
			}
		}
	}
	return nil
}

func (instance *ObjectHandler) decodeUnstructured(content []byte) (*unstructured.Unstructured, error) {
	result := &unstructured.Unstructured{}

	_, _, err := instance.Deserializer.Decode(content, nil, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (instance *ObjectHandler) newDecoder(content []byte) streaming.Decoder {
	buf := common.ToReadCloser(bytes.NewReader(content))
	fr := json.YAMLFramer.NewFrameReader(buf)
	return streaming.NewDecoder(fr, instance.Deserializer)
}
