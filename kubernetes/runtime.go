package kubernetes

import (
	"github.com/googleapis/gnostic/OpenAPIv2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	dynamicFake "k8s.io/client-go/dynamic/fake"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Runtime interface {
	ContextName() string
	NewDynamicClient() (dynamic.Interface, error)

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
	config      *restclient.Config
	contextName string

	discoveryClient *discovery.CachedDiscoveryClient
}

func (instance *runtimeImpl) NewDynamicClient() (dynamic.Interface, error) {
	return dynamic.NewForConfig(instance.config)
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

func (instance *runtimeMock) ContextName() string {
	return instance.contextName
}

func (instance *runtimeMock) OpenAPISchema() (*openapi_v2.Document, error) {
	return &openapi_v2.Document{}, nil
}
