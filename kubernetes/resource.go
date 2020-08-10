package kubernetes

import (
	"fmt"
	"github.com/echocat/kubor/kubernetes/transformation"
	"github.com/echocat/kubor/model"
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

func GetObjectResource(object *unstructured.Unstructured, client dynamic.Interface, by ObjectValidator) (ObjectResource, error) {
	info, err := GetObjectInfo(object, by)
	if err != nil {
		return ObjectResource{}, err
	}
	resource := client.Resource(info.Resource).Namespace(info.Namespace.String())
	return ObjectResource{
		ObjectInfo: info,
		Client:     client,
		Resource:   resource,
		Object:     object,
	}, nil
}

func (instance ObjectResource) Clone() ObjectResource {
	return ObjectResource{
		ObjectInfo: instance.ObjectInfo,
		Client:     instance.Client,
		Resource:   instance.Resource,
		Object:     instance.Object.DeepCopy(),
	}
}

func (instance ObjectResource) CloneForCreate(project *model.Project) (ObjectResource, error) {
	result := instance.Clone()
	if err := transformation.TransformForCreate(project, result.Object); err != nil {
		return ObjectResource{}, err
	}
	return result, nil
}

func (instance ObjectResource) CloneForUpdate(project *model.Project, original unstructured.Unstructured) (ObjectResource, error) {
	result := instance.Clone()
	if err := transformation.TransformForUpdate(project, original, result.Object); err != nil {
		return ObjectResource{}, err
	}
	return result, nil
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
	err := instance.Resource.Delete(instance.Name.String(), options, subresources...)
	return OptimizeError(err)
}

func (instance ObjectResource) Get(options *metav1.GetOptions, subresources ...string) (*unstructured.Unstructured, error) {
	if options == nil {
		options = &metav1.GetOptions{}
	}
	options.TypeMeta = instance.TypeMeta
	result, err := instance.Resource.Get(instance.Name.String(), *options, subresources...)
	return result, OptimizeError(err)
}

func (instance ObjectResource) Watch(options *metav1.ListOptions) (watch.Interface, error) {
	if options == nil {
		options = &metav1.ListOptions{}
	}
	options.TypeMeta = instance.TypeMeta
	options.Watch = true
	options.FieldSelector = fmt.Sprintf("metadata.name=%v", instance.Name)
	result, err := instance.Resource.Watch(*options)
	return result, OptimizeError(err)
}
