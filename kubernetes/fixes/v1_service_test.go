package fixes

import (
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func Test_updateFixForV1ServiceIfClusterIpIsAbsent_ignores_different_GroupVersionKinds(t *testing.T) {
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
	err := updateFixForV1ServiceIfClusterIpIsAbsent(mockedPih{}, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_updateFixForV1ServiceIfClusterIpIsAbsent_ignores_on_non_Service_kind(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Servicex",
		},
	}
	expectedTarget := original
	target := original
	err := updateFixForV1ServiceIfClusterIpIsAbsent(mockedPih{}, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_updateFixForV1ServiceIfClusterIpIsAbsent_ignores_on_non_v1_version(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1x",
			"kind":       "Service",
		},
	}
	expectedTarget := original
	target := original
	err := updateFixForV1ServiceIfClusterIpIsAbsent(mockedPih{}, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_updateFixForV1ServiceIfClusterIpIsAbsent_ignores_if_original_has_no_clusterIp(t *testing.T) {
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
	err := updateFixForV1ServiceIfClusterIpIsAbsent(mockedPih{}, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(expectedTarget).To(Equal(target))
}

func Test_updateFixForV1ServiceIfClusterIpIsAbsent_ignores_if_original_clusterIp_is_not_a_string(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"spec": map[string]interface{}{
				"clusterIP": 666,
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
	err := updateFixForV1ServiceIfClusterIpIsAbsent(mockedPih{}, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_updateFixForV1ServiceIfClusterIpIsAbsent_fails_if_target_has_spec_which_is_not_a_map(t *testing.T) {
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
	err := updateFixForV1ServiceIfClusterIpIsAbsent(mockedPih{}, original, &target)
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal("'spec' property of target does already exists but is not of type map[string]interface{} it is int"))
}

func Test_updateFixForV1ServiceIfClusterIpIsAbsent_set_spec_and_clusterIP(t *testing.T) {
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
	err := updateFixForV1ServiceIfClusterIpIsAbsent(mockedPih{}, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_updateFixForV1ServiceIfClusterIpIsAbsent_set_clusterIP(t *testing.T) {
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
	err := updateFixForV1ServiceIfClusterIpIsAbsent(mockedPih{}, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_updateFixForV1ServiceIfResourceVersionIsAbsent_ignores_different_GroupVersionKinds(t *testing.T) {
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
	err := updateFixForV1ServiceIfResourceVersionIsAbsent(mockedPih{}, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_updateFixForV1ServiceIfResourceVersionIsAbsent_ignores_on_non_Service_kind(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Servicex",
		},
	}
	expectedTarget := original
	target := original
	err := updateFixForV1ServiceIfResourceVersionIsAbsent(mockedPih{}, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_updateFixForV1ServiceIfResourceVersionIsAbsent_ignores_on_non_v1_version(t *testing.T) {
	g := NewGomegaWithT(t)

	original := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1x",
			"kind":       "Service",
		},
	}
	expectedTarget := original
	target := original
	err := updateFixForV1ServiceIfResourceVersionIsAbsent(mockedPih{}, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_updateFixForV1ServiceIfResourceVersionIsAbsent_ignores_if_original_has_no_resourceVersion(t *testing.T) {
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
	err := updateFixForV1ServiceIfResourceVersionIsAbsent(mockedPih{}, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(expectedTarget).To(Equal(target))
}

func Test_updateFixForV1ServiceIfResourceVersionIsAbsent_ignores_if_original_resourceVersion_is_not_a_string(t *testing.T) {
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
	err := updateFixForV1ServiceIfResourceVersionIsAbsent(mockedPih{}, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_updateFixForV1ServiceIfResourceVersionIsAbsent_fails_if_target_has_metadata_which_is_not_a_map(t *testing.T) {
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
	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata":   666,
		},
	}
	err := updateFixForV1ServiceIfResourceVersionIsAbsent(mockedPih{}, original, &target)
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal("'metadata' property of target does already exists but is not of type map[string]interface{} it is int"))
}

func Test_updateFixForV1ServiceIfResourceVersionIsAbsent_set_metadata_and_resourceVersion(t *testing.T) {
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
	err := updateFixForV1ServiceIfResourceVersionIsAbsent(mockedPih{}, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

func Test_updateFixForV1ServiceIfResourceVersionIsAbsent_set_resourceVersion(t *testing.T) {
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
	err := updateFixForV1ServiceIfResourceVersionIsAbsent(mockedPih{}, original, &target)
	g.Expect(err).To(BeNil())
	g.Expect(target).To(Equal(expectedTarget))
}

type mockedPih struct {
}

func (instance mockedPih) GetGroupId() string {
	panic("not implemented")
}

func (instance mockedPih) GetArtifactId() string {
	panic("not implemented")
}

func (instance mockedPih) GetRelease() string {
	panic("not implemented")
}
