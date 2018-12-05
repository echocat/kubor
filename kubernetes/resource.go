package kubernetes

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
)

type ObjectResource struct {
	ObjectInfo
	Client   dynamic.Interface
	Resource dynamic.ResourceInterface
	Object   *unstructured.Unstructured
}

func GetObjectResource(object *unstructured.Unstructured, client dynamic.Interface) (ObjectResource, error) {
	info, err := GetObjectInfo(object)
	if err != nil {
		return ObjectResource{}, err
	}
	resource := client.Resource(info.GroupVersionResource).Namespace(info.Namespace)
	return ObjectResource{
		ObjectInfo: info,
		Client:     client,
		Resource:   resource,
		Object:     object,
	}, nil
}

func (instance ObjectResource) Create(options *metav1.CreateOptions, subresources ...string) (*unstructured.Unstructured, error) {
	if options == nil {
		options = &metav1.CreateOptions{}
	}
	options.TypeMeta = instance.TypeMeta
	result, err := instance.Resource.Create(instance.Object, *options, subresources...)
	return result, OptimizeError(err)
}

func (instance ObjectResource) Update(options *metav1.UpdateOptions, subresources ...string) (*unstructured.Unstructured, error) {
	if options == nil {
		options = &metav1.UpdateOptions{}
	}
	options.TypeMeta = instance.TypeMeta
	result, err := instance.Resource.Update(instance.Object, *options, subresources...)
	return result, OptimizeError(err)
}

func (instance ObjectResource) Delete(options *metav1.DeleteOptions, subresources ...string) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	options.TypeMeta = instance.TypeMeta
	err := instance.Resource.Delete(instance.Name, options, subresources...)
	return OptimizeError(err)
}

func (instance ObjectResource) Get(options *metav1.GetOptions, subresources ...string) (*unstructured.Unstructured, error) {
	if options == nil {
		options = &metav1.GetOptions{}
	}
	options.TypeMeta = instance.TypeMeta
	result, err := instance.Resource.Get(instance.Name, *options, subresources...)
	return result, OptimizeError(err)
}

func (instance ObjectResource) Watch(options *metav1.ListOptions) (watch.Interface, error) {
	if options == nil {
		options = &metav1.ListOptions{}
	}
	options.TypeMeta = instance.TypeMeta
	options.Watch = true
	options.FieldSelector = fmt.Sprintf("metadata.name=%s", instance.Name)
	result, err := instance.Resource.Watch(*options)
	return result, OptimizeError(err)
}
