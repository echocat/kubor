package transformation

import (
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func Test_fixV1ServiceIfClusterIpIsAbsent_ignores_different_GroupVersionKinds(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
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
	err := fixV1ServiceIfClusterIpIsAbsent(nil, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_fixV1ServiceIfClusterIpIsAbsent_ignores_on_non_Service_kind(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Servicex",
		},
	}
	expectedTarget := original
	target := original
	err := fixV1ServiceIfClusterIpIsAbsent(nil, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_fixV1ServiceIfClusterIpIsAbsent_ignores_on_non_v1_version(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1x",
			"kind":       "Service",
		},
	}
	expectedTarget := original
	target := original
	err := fixV1ServiceIfClusterIpIsAbsent(nil, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_fixV1ServiceIfClusterIpIsAbsent_ignores_if_original_has_no_clusterIp(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
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
	err := fixV1ServiceIfClusterIpIsAbsent(nil, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(expectedTarget).To(Equal(target))
}

func Test_fixV1ServiceIfClusterIpIsAbsent_fails_if_target_has_spec_which_is_not_a_map(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"clusterIP": "1.2.3.4",
			},
		},
	}
	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec":       666,
		},
	}
	err := fixV1ServiceIfClusterIpIsAbsent(nil, original, &target)
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal(".spec.clusterIP accessor error: 666 is of the type int, expected map[string]interface{}"))
}

func Test_fixV1ServiceIfClusterIpIsAbsent_set_spec_and_clusterIP(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"clusterIP": "1.2.3.4",
			},
		},
	}
	expectedTarget := original
	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
		},
	}
	err := fixV1ServiceIfClusterIpIsAbsent(nil, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_fixV1ServiceIfClusterIpIsAbsent_set_clusterIP(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
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
	err := fixV1ServiceIfClusterIpIsAbsent(nil, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}
