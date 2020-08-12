package kubernetes

import (
	"errors"
	"fmt"
	"github.com/echocat/kubor/model"
	"github.com/googleapis/gnostic/OpenAPIv2"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"reflect"
)

func ResolveDryRun(in model.DryRunOn, gvk model.GroupVersionKind, client dynamic.Interface, runtime Runtime) (model.DryRunOn, error) {
	if in == model.DryRunOnServerIfPossible || in == model.DryRunOnServer {
		if serverSidePossible, err := HasServerDryRunSupport(gvk, client, runtime); err != nil {
			return "", err
		} else if in == model.DryRunOnServer {
			if !serverSidePossible {
				return "", fmt.Errorf("%v does not support server side dry run", gvk)
			}
		} else if serverSidePossible {
			return model.DryRunOnServer, nil
		} else {
			return model.DryRunOnClient, nil
		}
	}
	return in, nil
}

func HasServerDryRunSupport(gvk model.GroupVersionKind, client dynamic.Interface, runtime Runtime) (bool, error) {
	oapi, err := runtime.OpenAPISchema()
	if err != nil {
		return false, fmt.Errorf("failed to download openapi: %w", err)
	}
	supports, err := supportsServerDryRun(oapi, gvk)
	if err != nil {
		// We assume that we couldn't find the type, then check for namespace:
		supports, _ = supportsServerDryRun(oapi, model.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"})
		// If namespace supports dryRun, then we will support dryRun for CRDs only.
		if supports {
			if supports, err = hasCrd(gvk.GroupKind(), client); err != nil {
				return false, fmt.Errorf("failed to check CRD: %w", err)
			}
		}
	}
	return supports, nil
}

func hasGvkExtension(extensions []*openapi_v2.NamedAny, gvk model.GroupVersionKind) bool {
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
func supportsServerDryRun(doc *openapi_v2.Document, gvk model.GroupVersionKind) (bool, error) {
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
		return nil, fmt.Errorf("failed to list CRDs: %w", err)
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
