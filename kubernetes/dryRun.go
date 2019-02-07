package kubernetes

import (
	"errors"
	"fmt"
	"github.com/googleapis/gnostic/OpenAPIv2"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"reflect"
)

const (
	NowhereDryRun          DryRunOn = "nowhere"
	ClientDryRun           DryRunOn = "client"
	ServerDryRun           DryRunOn = "server"
	ServerIfPossibleDryRun DryRunOn = "serverIfPossible"
)

type DryRunOn string

func (instance DryRunOn) String() string {
	return string(instance)
}

func (instance *DryRunOn) Set(plain string) error {
	candidate := DryRunOn(plain)
	switch candidate {
	case NowhereDryRun, ClientDryRun, ServerDryRun, ServerIfPossibleDryRun:
		*instance = candidate
		return nil
	default:
		return fmt.Errorf("unknown dryRunOn: %v", candidate)
	}
}

func (instance DryRunOn) IsEnabled() bool {
	switch instance {
	case ServerIfPossibleDryRun:
		panic("This should not be called directly. Please resolve this before in your code to either server or client.")
	case ClientDryRun, ServerDryRun:
		return true
	default:
		return false
	}
}

func (instance DryRunOn) Resolve(gvk schema.GroupVersionKind, client dynamic.Interface, runtime Runtime) (DryRunOn, error) {
	if instance == ServerIfPossibleDryRun || instance == ServerDryRun {
		if serverSidePossible, err := HasServerDryRunSupport(gvk, client, runtime); err != nil {
			return DryRunOn(""), err
		} else if instance == ServerDryRun {
			if !serverSidePossible {
				return DryRunOn(""), fmt.Errorf("%v does not support server side dry run", gvk)
			}
		} else if serverSidePossible {
			return ServerDryRun, nil
		} else {
			return ClientDryRun, nil
		}
	}
	return instance, nil
}

func (instance *DryRunOn) Get() interface{} {
	return instance
}

func HasServerDryRunSupport(gvk schema.GroupVersionKind, client dynamic.Interface, runtime Runtime) (bool, error) {
	oapi, err := runtime.OpenAPISchema()
	if err != nil {
		return false, fmt.Errorf("failed to download openapi: %v", err)
	}
	supports, err := supportsServerDryRun(oapi, gvk)
	if err != nil {
		// We assume that we couldn't find the type, then check for namespace:
		supports, _ = supportsServerDryRun(oapi, schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"})
		// If namespace supports dryRun, then we will support dryRun for CRDs only.
		if supports {
			if supports, err = hasCrd(gvk.GroupKind(), client); err != nil {
				return false, fmt.Errorf("failed to check CRD: %v", err)
			}
		}
	}
	return supports, nil
}

func hasGvkExtension(extensions []*openapi_v2.NamedAny, gvk schema.GroupVersionKind) bool {
	for _, extension := range extensions {
		if extension.GetValue().GetYaml() == "" ||
			extension.GetName() != "x-kubernetes-group-version-kind" {
			continue
		}
		var value map[string]string
		err := yaml.Unmarshal([]byte(extension.GetValue().GetYaml()), &value)
		if err != nil {
			continue
		}

		if value["group"] == gvk.Group && value["kind"] == gvk.Kind && value["version"] == gvk.Version {
			return true
		}
		return false
	}
	return false
}

// SupportsDryRun is a method that let's us look in the OpenAPI if the
// specific group-version-kind supports the dryRun query parameter for
// the PATCH end-point.
func supportsServerDryRun(doc *openapi_v2.Document, gvk schema.GroupVersionKind) (bool, error) {
	for _, path := range doc.GetPaths().GetPath() {
		// Is this describing the gvk we're looking for?
		if !hasGvkExtension(path.GetValue().GetPatch().GetVendorExtension(), gvk) {
			continue
		}
		for _, param := range path.GetValue().GetPatch().GetParameters() {
			if param.GetParameter().GetNonBodyParameter().GetQueryParameterSubSchema().GetName() == "dryRun" {
				return true, nil
			}
		}
		return false, nil
	}

	return false, errors.New("couldn't find GVK in openapi")
}

func crdFromDynamic(client dynamic.Interface) ([]schema.GroupKind, error) {
	list, err := client.Resource(schema.GroupVersionResource{
		Group:    "apiextensions.k8s.io",
		Version:  "v1beta1",
		Resource: "customresourcedefinitions",
	}).List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list CRDs: %v", err)
	}
	if list == nil {
		return nil, nil
	}

	var gks []schema.GroupKind

	// We need to parse the list to get the gvk, I guess that's fine.
	for _, crd := range (*list).Items {
		// Look for group, version, and kind
		group, _, _ := unstructured.NestedString(crd.Object, "spec", "group")
		kind, _, _ := unstructured.NestedString(crd.Object, "spec", "names", "kind")

		gks = append(gks, schema.GroupKind{
			Group: group,
			Kind:  kind,
		})
	}

	return gks, nil
}

func hasCrd(gvk schema.GroupKind, client dynamic.Interface) (bool, error) {
	list, err := crdFromDynamic(client)
	if err != nil {
		return false, err
	}

	for _, crd := range list {
		if reflect.DeepEqual(gvk, crd) {
			return true, nil
		}
	}
	return false, nil
}
