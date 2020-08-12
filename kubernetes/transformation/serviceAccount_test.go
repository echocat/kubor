package transformation

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func Test_preserveServiceAccountSecrets_ignores_different_GroupVersionKinds(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1beta1",
			"kind":       "ServiceAccount",
		},
	}
	target := expectedTarget

	err := preserveServiceAccountSecrets(nil, existing, &target, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceAccountSecrets_ignores_on_non_ServiceAccount_kind(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccountx",
		},
	}
	expectedTarget := existing
	target := existing

	err := preserveServiceAccountSecrets(nil, existing, &target, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceAccountSecrets_ignores_if_target_has_already_secrets(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"secrets": []interface{}{
				map[string]interface{}{
					"kind":      "anExistingKind",
					"namespace": "anExistingNamespace",
					"name":      "anExistingName",
				},
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"secrets": []interface{}{
				map[string]interface{}{
					"kind":      "aTargetKind",
					"namespace": "aTargetNamespace",
					"name":      "aTargetName",
				},
			},
		},
	}
	target := expectedTarget

	err := preserveServiceAccountSecrets(nil, existing, &target, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceAccountSecrets_ignores_if_target_has_already_imagePullSecrets(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"imagePullSecrets": []interface{}{
				map[string]interface{}{
					"name": "anExistingName",
				},
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"imagePullSecrets": []interface{}{
				map[string]interface{}{
					"name": "aTargetName",
				},
			},
		},
	}
	target := expectedTarget

	err := preserveServiceAccountSecrets(nil, existing, &target, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceAccountSecrets_applies_secrets(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"secrets": []interface{}{
				map[string]interface{}{
					"kind":      "anExistingKind",
					"namespace": "anExistingNamespace",
					"name":      "anExistingName",
				},
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"secrets": []interface{}{
				map[string]interface{}{
					"kind":      "anExistingKind",
					"namespace": "anExistingNamespace",
					"name":      "anExistingName",
				},
			},
			"imagePullSecrets": []interface{}(nil),
		},
	}
	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
		},
	}

	err := preserveServiceAccountSecrets(nil, existing, &target, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceAccountSecrets_applies_imagePullSecrets(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"imagePullSecrets": []interface{}{
				map[string]interface{}{
					"name": "anExistingName",
				},
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"secrets":    []interface{}(nil),
			"imagePullSecrets": []interface{}{
				map[string]interface{}{
					"name": "anExistingName",
				},
			},
		},
	}
	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
		},
	}

	err := preserveServiceAccountSecrets(nil, existing, &target, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}
