package kubernetes

import (
	openapi_v2 "github.com/google/gnostic/openapiv2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	dynamicFake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/rest"
	restFake "k8s.io/client-go/rest/fake"
	"k8s.io/client-go/tools/clientcmd"
)

type Runtime interface {
	ContextName() string
	NewDynamicClient() (dynamic.Interface, error)
	NewRestClient(gvk schema.GroupVersionKind) (rest.Interface, error)

	discovery.OpenAPISchemaInterface
}

func newRuntimeImpl(clientConfig clientcmd.ClientConfig, contextName string) (*runtimeImpl, error) {
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	dc, err := newDiscoveryClientFor(config)
	if err != nil {
		return nil, err
	}
	return &runtimeImpl{
		config:          config,
		contextName:     contextName,
		discoveryClient: dc,
	}, nil
}

type runtimeImpl struct {
	config      *rest.Config
	contextName string

	discoveryClient discovery.DiscoveryInterface
}

func (instance *runtimeImpl) NewDynamicClient() (dynamic.Interface, error) {
	return dynamic.NewForConfig(instance.config)
}

func (instance *runtimeImpl) NewRestClient(gvk schema.GroupVersionKind) (rest.Interface, error) {
	config := dynamic.ConfigFor(instance.config)
	gv := gvk.GroupVersion()
	config.GroupVersion = &gv
	config.APIPath = "/api"
	return rest.RESTClientFor(config)
}

func (instance *runtimeImpl) ContextName() string {
	return instance.contextName
}

func (instance *runtimeImpl) OpenAPISchema() (*openapi_v2.Document, error) {
	return instance.discoveryClient.OpenAPISchema()
}

func newRuntimeMock(contextName string) (*runtimeMock, error) {
	return &runtimeMock{
		scheme:      runtime.NewScheme(),
		contextName: contextName,
	}, nil
}

type runtimeMock struct {
	scheme      *runtime.Scheme
	contextName string
}

func (instance *runtimeMock) NewDynamicClient() (dynamic.Interface, error) {
	return dynamicFake.NewSimpleDynamicClient(instance.scheme), nil
}

func (instance *runtimeMock) NewRestClient(gvk schema.GroupVersionKind) (rest.Interface, error) {
	return &restFake.RESTClient{
		GroupVersion: gvk.GroupVersion(),
	}, nil
}

func (instance *runtimeMock) ContextName() string {
	return instance.contextName
}

func (instance *runtimeMock) OpenAPISchema() (*openapi_v2.Document, error) {
	return &openapi_v2.Document{}, nil
}
