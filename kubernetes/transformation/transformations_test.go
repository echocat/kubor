package transformation

import (
	"fmt"
	"github.com/echocat/kubor/model"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

var testTransformations = Transformations{
	Updates: map[Name]Update{
		"ufoo": testTransformation{"ufoo"},
		"ubar": testTransformation{"ubar"},
		"foo":  testTransformation{"foo"},
		"bar":  testTransformation{"bar"},
	},
	Creates: map[Name]Create{
		"cfoo": testTransformation{"cfoo"},
		"cbar": testTransformation{"cbar"},
		"foo":  testTransformation{"foo"},
		"bar":  testTransformation{"bar"},
	},
}

func TestTransformations_TransformForUpdate_calling_enabled(t *testing.T) {
	project := model.NewProject()

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
					model.AnnotationTransformationPrefix + "foo":  "something else",
					model.AnnotationTransformationPrefix + "bar":  "disabled",
					model.AnnotationTransformationPrefix + "ufoo": "enabled",
					model.AnnotationTransformationPrefix + "ubar": "disabled",
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
					model.AnnotationTransformationPrefix + "foo":  "something else",
					model.AnnotationTransformationPrefix + "bar":  "disabled",
					model.AnnotationTransformationPrefix + "ufoo": "enabled",
					model.AnnotationTransformationPrefix + "ubar": "disabled",
					"foo":  "test:foo:something else",
					"ufoo": "test:ufoo:",
				},
			},
		},
	}

	err := testTransformations.TransformForUpdate(&project, existing, &target)

	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func TestTransformations_TransformForCreate_calling_enabled(t *testing.T) {
	project := model.NewProject()

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
					model.AnnotationTransformationPrefix + "foo":  "something else",
					model.AnnotationTransformationPrefix + "bar":  "disabled",
					model.AnnotationTransformationPrefix + "cfoo": "enabled",
					model.AnnotationTransformationPrefix + "cbar": "disabled",
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
					model.AnnotationTransformationPrefix + "foo":  "something else",
					model.AnnotationTransformationPrefix + "bar":  "disabled",
					model.AnnotationTransformationPrefix + "cfoo": "enabled",
					model.AnnotationTransformationPrefix + "cbar": "disabled",
					"foo":  "test:foo:something else",
					"cfoo": "test:cfoo:",
				},
			},
		},
	}

	err := testTransformations.TransformForCreate(&project, &target)

	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

type testTransformation struct {
	name Name
}

func (instance testTransformation) TransformForUpdate(p *model.Project, _ unstructured.Unstructured, target *unstructured.Unstructured, argument string) error {
	return instance.TransformForCreate(p, target, argument)
}

func (instance testTransformation) TransformForCreate(_ *model.Project, target *unstructured.Unstructured, argument string) error {
	if argument == "error" {
		return fmt.Errorf("expected:%v:%s", instance.name, argument)
	}

	as := target.GetAnnotations()
	as[instance.name.String()] = fmt.Sprintf("test:%v:%s", instance.name, argument)
	target.SetAnnotations(as)

	return nil
}