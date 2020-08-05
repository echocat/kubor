package transformation

import (
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func Test_fixIfResourceVersionIsAbsent_ignores_different_GroupVersionKinds(t *testing.T) {
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
	err := fixIfResourceVersionIsAbsent(nil, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_fixIfResourceVersionIsAbsent_ignores_on_non_Service_kind(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Servicex",
		},
	}
	expectedTarget := original
	target := original
	err := fixIfResourceVersionIsAbsent(nil, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_fixIfResourceVersionIsAbsent_ignores_on_non_v1_version(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1x",
			"kind":       "Service",
		},
	}
	expectedTarget := original
	target := original
	err := fixIfResourceVersionIsAbsent(nil, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_fixIfResourceVersionIsAbsent_ignores_if_original_has_no_resourceVersion(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
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
	err := fixIfResourceVersionIsAbsent(nil, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(expectedTarget).To(Equal(target))
}

func Test_fixIfResourceVersionIsAbsent_ignores_if_original_resourceVersion_is_not_a_string(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
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
	err := fixIfResourceVersionIsAbsent(nil, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_fixIfResourceVersionIsAbsent_set_metadata_and_resourceVersion(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"resourceVersion": "1.2.3.4",
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
	err := fixIfResourceVersionIsAbsent(nil, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_fixIfResourceVersionIsAbsent_set_resourceVersion(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
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
	err := fixIfResourceVersionIsAbsent(nil, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}
