package model

import (
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"sync"
)

var DefaultGroupVersionKindRegistry = (&GroupVersionKindRegistry{}).
	With(BuildGroupVersionKinds(rbacv1.SchemeGroupVersion, &rbacv1.Role{}).Build()).
	With(BuildGroupVersionKinds(rbacv1.SchemeGroupVersion, &rbacv1.RoleBinding{}).Build()).
	With(BuildGroupVersionKinds(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRole{}).Build()).
	With(BuildGroupVersionKinds(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRoleBinding{}).Build()).
	With(BuildGroupVersionKinds(v1.SchemeGroupVersion, &v1.Namespace{}).Build()).
	With(BuildGroupVersionKinds(v1.SchemeGroupVersion, &v1.Service{}).Build()).
	With(BuildGroupVersionKinds(v1.SchemeGroupVersion, &v1.Secret{}).Build()).
	With(BuildGroupVersionKinds(v1.SchemeGroupVersion, &v1.ServiceAccount{}).Build()).
	With(BuildGroupVersionKinds(v1.SchemeGroupVersion, &v1.Pod{}).Build()).
	With(BuildGroupVersionKinds(v1.SchemeGroupVersion, &v1.PersistentVolume{}).Build()).
	With(BuildGroupVersionKinds(v1.SchemeGroupVersion, &v1.PersistentVolumeClaim{}).Build()).
	With(BuildGroupVersionKinds(v1.SchemeGroupVersion, &v1.ConfigMap{}).Build()).
	With(BuildGroupVersionKinds(appsv1.SchemeGroupVersion, &appsv1.Deployment{}).
		With(appsv1beta1.SchemeGroupVersion, &appsv1beta1.Deployment{}).
		With(appsv1beta2.SchemeGroupVersion, &appsv1beta2.Deployment{}).
		With(extensionsv1beta1.SchemeGroupVersion, &extensionsv1beta1.Deployment{}).
		Build()).
	With(BuildGroupVersionKinds(appsv1.SchemeGroupVersion, &appsv1.StatefulSet{}).
		With(appsv1beta1.SchemeGroupVersion, &appsv1beta1.StatefulSet{}).
		With(appsv1beta2.SchemeGroupVersion, &appsv1beta2.StatefulSet{}).
		Build()).
	With(BuildGroupVersionKinds(appsv1.SchemeGroupVersion, &appsv1.DaemonSet{}).
		With(appsv1beta2.SchemeGroupVersion, &appsv1beta2.DaemonSet{}).
		With(extensionsv1beta1.SchemeGroupVersion, &extensionsv1beta1.DaemonSet{}).
		Build()).
	With(BuildGroupVersionKinds(appsv1.SchemeGroupVersion, &appsv1.ReplicaSet{}).
		With(appsv1beta2.SchemeGroupVersion, &appsv1beta2.ReplicaSet{}).
		With(extensionsv1beta1.SchemeGroupVersion, &extensionsv1beta1.ReplicaSet{}).
		Build()).
	With(BuildGroupVersionKinds(apiextensionsv1.SchemeGroupVersion, &apiextensionsv1.CustomResourceDefinition{}).
		With(apiextensions.SchemeGroupVersion, &apiextensions.CustomResourceDefinition{}).
		With(apiextensionsv1beta1.SchemeGroupVersion, &apiextensionsv1beta1.CustomResourceDefinition{}).
		Build()).
	With(BuildGroupVersionKinds(networkingv1.SchemeGroupVersion, &networkingv1.NetworkPolicy{}).
		With(extensionsv1beta1.SchemeGroupVersion, &extensionsv1beta1.NetworkPolicy{}).
		Build()).
	With(BuildGroupVersionKinds(networkingv1beta1.SchemeGroupVersion, &networkingv1beta1.Ingress{}).
		With(extensionsv1beta1.SchemeGroupVersion, &extensionsv1beta1.Ingress{}).
		Build()).
	With(BuildGroupVersionKinds(batchv1.SchemeGroupVersion, &batchv1.Job{}).Build()).
	With(BuildGroupVersionKinds(batchv1beta1.SchemeGroupVersion, &batchv1beta1.CronJob{}).
		With(batchv1beta1.SchemeGroupVersion, &batchv1beta1.CronJob{}).
		Build())

type GroupVersionKindRegistry struct {
	assignments map[GroupVersionKind]*GroupVersionKinds
	mutex       sync.RWMutex
}

func (instance *GroupVersionKindRegistry) With(twins GroupVersionKinds) *GroupVersionKindRegistry {
	if instance == nil {
		*instance = GroupVersionKindRegistry{}
	}

	instance.mutex.Lock()
	defer instance.mutex.Unlock()

	if instance.assignments == nil {
		instance.assignments = make(map[GroupVersionKind]*GroupVersionKinds)
	}

	for actual := range twins {
		if existing := instance.assignments[actual]; existing != nil {
			delete(*existing, actual)
		}
		instance.assignments[actual] = &twins
	}

	return instance
}

func (instance *GroupVersionKindRegistry) AreTwins(left, right GroupVersionKind) bool {
	if instance == nil {
		return false
	}

	instance.mutex.RLock()
	defer instance.mutex.RUnlock()

	if instance.assignments == nil {
		return false
	}

	existing := instance.assignments[left]

	if existing == nil {
		return false
	}

	return existing.Contains(right)
}

func (instance *GroupVersionKindRegistry) GetTwins(of GroupVersionKind) (result GroupVersionKinds) {
	result = GroupVersionKinds{}
	if instance == nil {
		return
	}

	instance.mutex.RLock()
	defer instance.mutex.RUnlock()

	if instance.assignments == nil {
		return
	}

	existing := instance.assignments[of]

	if existing != nil && *existing != nil {
		for candidate := range *existing {
			if candidate != of {
				result[candidate] = true
			}
		}
	}

	return
}

func (instance *GroupVersionKindRegistry) AsGroupVersionKinds() (result GroupVersionKinds) {
	result = GroupVersionKinds{}
	if instance == nil {
		return
	}

	instance.mutex.RLock()
	defer instance.mutex.RUnlock()

	if instance.assignments != nil {
		for candidate := range instance.assignments {
			result[candidate] = true
		}
	}

	return
}
