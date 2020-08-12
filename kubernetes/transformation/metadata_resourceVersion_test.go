package transformation

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func Test_preserveResourceVersion_ignores_different_GroupVersionKinds(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Service",
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "extensions/v1",
			"kind":       "Service",
		},
	}
	target := expectedTarget

	err := preserveResourceVersion(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveResourceVersion_ignores_on_non_Service_kind(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Servicex",
		},
	}
	expectedTarget := existing
	target := existing

	err := preserveResourceVersion(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveResourceVersion_ignores_on_non_v1_version(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1x",
			"kind":       "Service",
		},
	}
	expectedTarget := existing
	target := existing

	err := preserveResourceVersion(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveResourceVersion_ignores_if_original_has_no_resourceVersion(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"foo": "bar",
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
		},
	}
	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
		},
	}

	err := preserveResourceVersion(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveResourceVersion_ignores_if_original_resourceVersion_is_not_a_string(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"resourceVersion": 666,
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
		},
	}
	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
		},
	}

	err := preserveResourceVersion(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveResourceVersion_set_metadata_and_resourceVersion(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"resourceVersion": "1.2.3.4",
			},
		},
	}
	expectedTarget := existing
	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
		},
	}

	err := preserveResourceVersion(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveResourceVersion_set_resourceVersion(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"resourceVersion": "1.2.3.4",
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"foo":             "bar",
				"resourceVersion": "1.2.3.4",
			},
		},
	}
	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"foo": "bar",
			},
		},
	}

	err := preserveResourceVersion(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}
