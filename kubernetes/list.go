package kubernetes

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"reflect"
)

type QueryNamespaceReceiver func(Object) error

func QueryNamespace(client dynamic.Interface, gvk schema.GroupVersionKind, namespace string, receiver QueryNamespaceReceiver) error {
	gvr, _ := meta.UnsafeGuessKindToResource(gvk)
	resource := client.Resource(gvr).Namespace(namespace)

	opts := metav1.ListOptions{
		Limit: 10,
	}
	for {
		list, err := resource.List(opts)
		if err != nil {
			return fmt.Errorf("cannot query namespace %s for %v: %v", namespace, gvk, err)
		}

		err = list.EachListItem(func(object runtime.Object) error {
			if ko, ok := object.(Object); ok {
				return receiver(ko)
			}
			return fmt.Errorf("object of type %v are not supported", reflect.TypeOf(object))
		})
		if err != nil {
			return fmt.Errorf("cannot query namespace %s for %v: %v", namespace, gvk, err)
		}

		opts.Continue = list.GetContinue()
		if opts.Continue == "" {
			return nil
		}
	}
}
