package transformation

import (
	"github.com/echocat/kubor/model"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func Test_ensureKuborAnnotations_drop(t *testing.T) {
	project := model.NewProject()
	project.Annotations.Transformations.Action = model.AnnotationActionDrop
	project.Annotations.Transformations.Name = "foo-"

	target := unstructured.Unstructured{
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
			"metadata": map[string]interface{}{
				"annotations": map[string]interface{}{
					"foo":     "bar",
					"foo-foo": "aFoo",
					"foo-bar": "aBar",
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
			"metadata": map[string]interface{}{
				"annotations": map[string]interface{}{
					"foo": "bar",
				},
			},
		},
	}

	err := ensureKuborAnnotations(&project, &target, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_ensureKuborAnnotations_leave(t *testing.T) {
	project := model.NewProject()
	project.Annotations.Transformations.Action = model.AnnotationActionLeave
	project.Annotations.Transformations.Name = "foo-"

	target := unstructured.Unstructured{
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
			"metadata": map[string]interface{}{
				"annotations": map[string]interface{}{
					"foo":     "bar",
					"foo-foo": "aFoo",
					"foo-bar": "aBar",
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
			"metadata": map[string]interface{}{
				"annotations": map[string]interface{}{
					"foo":     "bar",
					"foo-foo": "aFoo",
					"foo-bar": "aBar",
				},
			},
		},
	}

	err := ensureKuborAnnotations(&project, &target, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_ensureKuborAnnotations_drop_with_spec_template(t *testing.T) {
	project := model.NewProject()
	project.Annotations.Transformations.Action = model.AnnotationActionDrop
	project.Annotations.Transformations.Name = "foo-"

	target := unstructured.Unstructured{
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
			"metadata": map[string]interface{}{
				"annotations": map[string]interface{}{
					"foo":     "bar",
					"foo-foo": "aFoo",
					"foo-bar": "aBar",
				},
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							"foo":     "bar2",
							"foo-foo": "aFoo2",
							"foo-bar": "aBar2",
						},
					},
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
			"metadata": map[string]interface{}{
				"annotations": map[string]interface{}{
					"foo": "bar",
				},
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							"foo": "bar2",
						},
					},
				},
			},
		},
	}

	err := ensureKuborAnnotations(&project, &target, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}
