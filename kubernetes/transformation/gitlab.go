package transformation

import (
	"github.com/echocat/kubor/model"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	gitlabDiscoveryTransformationName = model.TransformationName("gitlab-discovery")
	gitlabBuildTransformationName     = model.TransformationName("gitlab-build")

	gitlabAnnotationApp        = "app.gitlab.com/app"
	gitlabAnnotationEnv        = "app.gitlab.com/env"
	gitlabAnnotationJobId      = "app.gitlab.com/job-id"
	gitlabAnnotationPipelineId = "app.gitlab.com/pipeline-id"
	gitlabAnnotationRunnerId   = "app.gitlab.com/runner-id"
	gitlabAnnotationProjectUrl = "app.gitlab.com/project-url"

	gitlabEnvProjectPathSlug = "CI_PROJECT_PATH_SLUG"
	gitlabEnvEnvironmentSlug = "CI_ENVIRONMENT_SLUG"
	gitlabEnvJobId           = "CI_JOB_ID"
	gitlabEnvPipelineId      = "CI_PIPELINE_ID"
	gitlabEnvRunnerId        = "CI_RUNNER_ID"
	gitlabEnvProjectUrl      = "CI_PROJECT_URL"
)

var GitlabDiscoveryReceiverGvks = model.BuildGroupVersionKinds(v1.SchemeGroupVersion, &v1.Pod{}).
	With(appsv1.SchemeGroupVersion, &appsv1.Deployment{}).
	With(appsv1beta1.SchemeGroupVersion, &appsv1beta1.Deployment{}).
	With(appsv1beta2.SchemeGroupVersion, &appsv1beta2.Deployment{}).
	With(extensionsv1beta1.SchemeGroupVersion, &extensionsv1beta1.Deployment{}).
	With(appsv1.SchemeGroupVersion, &appsv1.StatefulSet{}).
	With(appsv1beta1.SchemeGroupVersion, &appsv1beta1.StatefulSet{}).
	With(appsv1beta2.SchemeGroupVersion, &appsv1beta2.StatefulSet{}).
	With(appsv1.SchemeGroupVersion, &appsv1.DaemonSet{}).
	With(appsv1beta2.SchemeGroupVersion, &appsv1beta2.DaemonSet{}).
	With(extensionsv1beta1.SchemeGroupVersion, &extensionsv1beta1.DaemonSet{}).
	With(appsv1.SchemeGroupVersion, &appsv1.ReplicaSet{}).
	With(appsv1beta2.SchemeGroupVersion, &appsv1beta2.ReplicaSet{}).
	With(extensionsv1beta1.SchemeGroupVersion, &extensionsv1beta1.ReplicaSet{}).
	With(batchv1.SchemeGroupVersion, &batchv1.Job{}).
	With(batchv1beta1.SchemeGroupVersion, &batchv1beta1.CronJob{}).
	Build()

func init() {
	Default.MustRegisterUpdateFunc(gitlabDiscoveryTransformationName, appendGitlabDiscoveryOnUpdate)
	Default.MustRegisterCreateFunc(gitlabDiscoveryTransformationName, appendGitlabDiscovery)

	Default.MustRegisterUpdateFunc(gitlabBuildTransformationName, appendGitlabBuildOnUpdate)
	Default.MustRegisterCreateFunc(gitlabBuildTransformationName, appendGitlabBuild)
}

func appendGitlabDiscoveryOnUpdate(project *model.Project, _ unstructured.Unstructured, target *unstructured.Unstructured, argument *string) error {
	return appendGitlabDiscovery(project, target, argument)
}

func appendGitlabDiscovery(project *model.Project, target *unstructured.Unstructured, _ *string) error {
	if !GitlabDiscoveryReceiverGvks.Contains(model.GroupVersionKind(target.GroupVersionKind())) {
		return nil
	}

	if err := appendGitlabDiscoveryOfPath(project, target, "metadata", "annotations"); err != nil {
		return err
	}

	if _, specTemplateExists, err := unstructured.NestedMap(target.Object, "spec", "template"); err != nil || !specTemplateExists {
		return err
	}
	if err := appendGitlabDiscoveryOfPath(project, target, "spec", "template", "metadata", "annotations"); err != nil {
		return err
	}

	return nil
}

func appendGitlabDiscoveryOfPath(project *model.Project, target *unstructured.Unstructured, fields ...string) error {
	annotations, _, err := NestedStringMap(target.Object, fields...)
	if err != nil {
		return err
	}
	if annotations == nil {
		annotations = make(map[string]string)
	}

	atLeastOneSet := false
	if v := project.Env[gitlabEnvProjectPathSlug]; v != "" {
		annotations[gitlabAnnotationApp] = v
		atLeastOneSet = true
	}
	if v := project.Env[gitlabEnvEnvironmentSlug]; v != "" {
		annotations[gitlabAnnotationEnv] = v
		atLeastOneSet = true
	}

	if !atLeastOneSet {
		return nil
	}

	return unstructured.SetNestedStringMap(target.Object, annotations, fields...)
}
func appendGitlabBuildOnUpdate(project *model.Project, _ unstructured.Unstructured, target *unstructured.Unstructured, argument *string) error {
	return appendGitlabBuild(project, target, argument)
}

func appendGitlabBuild(project *model.Project, target *unstructured.Unstructured, _ *string) error {
	if err := appendGitlabBuildOfPath(project, target, "metadata", "annotations"); err != nil {
		return err
	}

	if _, specTemplateExists, err := NestedMap(target.Object, "spec", "template"); err != nil || !specTemplateExists {
		return err
	}
	if err := appendGitlabBuildOfPath(project, target, "spec", "template", "metadata", "annotations"); err != nil {
		return err
	}

	return nil
}

func appendGitlabBuildOfPath(project *model.Project, target *unstructured.Unstructured, fields ...string) error {
	annotations, _, err := NestedStringMap(target.Object, fields...)
	if err != nil {
		return err
	}
	if annotations == nil {
		annotations = make(map[string]string)
	}

	atLeastOneSet := false
	if v := project.Env[gitlabEnvJobId]; v != "" {
		annotations[gitlabAnnotationJobId] = v
		atLeastOneSet = true
	}
	if v := project.Env[gitlabEnvPipelineId]; v != "" {
		annotations[gitlabAnnotationPipelineId] = v
		atLeastOneSet = true
	}
	if v := project.Env[gitlabEnvRunnerId]; v != "" {
		annotations[gitlabAnnotationRunnerId] = v
		atLeastOneSet = true
	}
	if v := project.Env[gitlabEnvProjectUrl]; v != "" {
		annotations[gitlabAnnotationProjectUrl] = v
		atLeastOneSet = true
	}

	if !atLeastOneSet {
		return nil
	}

	return unstructured.SetNestedStringMap(target.Object, annotations, fields...)
}
