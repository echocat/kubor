package transformation

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func Test_preserveServiceClusterIps_ignores_different_GroupVersionKinds(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"clusterIP": "1.2.3.4",
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "extensions/v1",
			"kind":       "Service",
		},
	}
	target := expectedTarget

	err := preserveServiceClusterIps(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceClusterIps_ignores_on_non_Service_kind(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Servicex",
			"spec": map[string]interface{}{
				"clusterIP": "1.2.3.4",
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Servicex",
		},
	}
	target := expectedTarget

	err := preserveServiceClusterIps(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceClusterIps_ignores_if_original_has_no_clusterIp(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
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

	err := preserveServiceClusterIps(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceClusterIps_set_spec_and_clusterIP(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"clusterIP": "1.2.3.4",
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

	err := preserveServiceClusterIps(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceClusterIps_does_not_override(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"clusterIP": "1.2.3.4",
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"clusterIP": "6.6.6.6",
			},
		},
	}
	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"clusterIP": "6.6.6.6",
			},
		},
	}

	err := preserveServiceClusterIps(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceClusterIps_set_clusterIP(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"clusterIP": "1.2.3.4",
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"foo":       "bar",
				"clusterIP": "1.2.3.4",
			},
		},
	}
	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"foo": "bar",
			},
		},
	}

	err := preserveServiceClusterIps(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceHealthCheckNodePort_ignores_different_GroupVersionKinds(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"healthCheckNodePort": "123",
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "extensions/v1",
			"kind":       "Service",
		},
	}
	target := expectedTarget

	err := preserveServiceHealthCheckNodePort(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceHealthCheckNodePort_ignores_on_non_Service_kind(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Servicex",
			"spec": map[string]interface{}{
				"healthCheckNodePort": "123",
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Servicex",
		},
	}
	target := expectedTarget

	err := preserveServiceHealthCheckNodePort(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceHealthCheckNodePort_ignores_if_original_has_no_clusterIp(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
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

	err := preserveServiceHealthCheckNodePort(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceHealthCheckNodePort_set_spec_and_clusterIP(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"healthCheckNodePort": "123",
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

	err := preserveServiceHealthCheckNodePort(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceHealthCheckNodePort_does_not_override(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"healthCheckNodePort": "123",
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"healthCheckNodePort": "666",
			},
		},
	}
	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"healthCheckNodePort": "666",
			},
		},
	}

	err := preserveServiceHealthCheckNodePort(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceHealthCheckNodePort_set_clusterIP(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"healthCheckNodePort": "123",
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"foo":                 "bar",
				"healthCheckNodePort": "123",
			},
		},
	}
	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"foo": "bar",
			},
		},
	}

	err := preserveServiceHealthCheckNodePort(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceNodePorts_ignores_different_GroupVersionKinds(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"ports": []interface{}{
					map[string]interface{}{
						"name":     "foo",
						"port":     int64(111),
						"nodePort": int64(123),
					},
					map[string]interface{}{
						"name":     "bar",
						"port":     int64(222),
						"nodePort": int64(234),
					},
				},
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "extensions/v1",
			"kind":       "Service",
		},
	}
	target := expectedTarget

	err := preserveServiceNodePorts(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceNodePorts_ignores_on_non_Service_kind(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Servicex",
			"spec": map[string]interface{}{
				"ports": []interface{}{
					map[string]interface{}{
						"name":     "foo",
						"port":     int64(111),
						"nodePort": int64(123),
					},
					map[string]interface{}{
						"name":     "bar",
						"port":     int64(222),
						"nodePort": int64(234),
					},
				},
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Servicex",
		},
	}
	target := expectedTarget

	err := preserveServiceNodePorts(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceNodePorts_ignores_if_original_has_no_nodePort(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"ports": []interface{}{
					map[string]interface{}{
						"name": "foo",
						"port": int64(111),
					},
					map[string]interface{}{
						"name": "bar",
						"port": int64(222),
					},
				},
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

	err := preserveServiceNodePorts(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_preserveServiceNodePorts_set(t *testing.T) {
	existing := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"ports": []interface{}{
					map[string]interface{}{
						"name":     "foo",
						"port":     int64(111),
						"nodePort": int64(123),
					},
					map[string]interface{}{
						"name":     "bar",
						"port":     int64(222),
						"nodePort": int64(234),
					},
				},
			},
		},
	}
	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"ports": []interface{}{
					map[string]interface{}{
						"name":     "foo",
						"port":     int64(111),
						"nodePort": int64(123),
					},
					map[string]interface{}{
						"name":     "bar",
						"port":     int64(333),
						"nodePort": int64(234),
					},
				},
			},
		},
	}
	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"ports": []interface{}{
					map[string]interface{}{
						"name": "foo",
						"port": int64(111),
					},
					map[string]interface{}{
						"name":     "bar",
						"port":     int64(333),
						"nodePort": int64(234),
					},
				},
			},
		},
	}

	err := preserveServiceNodePorts(nil, existing, &target, "")
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}
