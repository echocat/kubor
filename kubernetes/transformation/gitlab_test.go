package transformation

import (
	"github.com/echocat/kubor/model"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

func Test_appendGitlabDiscovery_append(t *testing.T) {
	project := model.NewProject()
	project.Env[gitlabEnvProjectPathSlug] = "/project-slug"
	project.Env[gitlabEnvEnvironmentSlug] = "/environment-slug"

	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
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

	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"annotations": map[string]interface{}{
					"foo":               "bar",
					gitlabAnnotationApp: "/project-slug",
					gitlabAnnotationEnv: "/environment-slug",
				},
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							"foo":               "bar2",
							gitlabAnnotationApp: "/project-slug",
							gitlabAnnotationEnv: "/environment-slug",
						},
					},
				},
			},
		},
	}

	err := appendGitlabDiscovery(&project, &target, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_appendGitlabDiscovery_does_not_append_if_env_absent(t *testing.T) {
	project := model.NewProject()

	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
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

	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
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

	err := appendGitlabDiscovery(&project, &target, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_appendGitlabBuild_append(t *testing.T) {
	project := model.NewProject()
	project.Env[gitlabEnvJobId] = "job1"
	project.Env[gitlabEnvPipelineId] = "pipeline2"
	project.Env[gitlabEnvRunnerId] = "runner3"
	project.Env[gitlabEnvProjectUrl] = "https://foo.bar"

	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
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

	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"annotations": map[string]interface{}{
					"foo":                      "bar",
					gitlabAnnotationJobId:      "job1",
					gitlabAnnotationPipelineId: "pipeline2",
					gitlabAnnotationRunnerId:   "runner3",
					gitlabAnnotationProjectUrl: "https://foo.bar",
				},
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"annotations": map[string]interface{}{
							"foo":                      "bar2",
							gitlabAnnotationJobId:      "job1",
							gitlabAnnotationPipelineId: "pipeline2",
							gitlabAnnotationRunnerId:   "runner3",
							gitlabAnnotationProjectUrl: "https://foo.bar",
						},
					},
				},
			},
		},
	}

	err := appendGitlabBuild(&project, &target, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}

func Test_appendGitlabBuild_does_not_append_if_env_absent(t *testing.T) {
	project := model.NewProject()

	target := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
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

	expectedTarget := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
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

	err := appendGitlabBuild(&project, &target, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
}
