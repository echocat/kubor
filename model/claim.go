package model

import (
	"fmt"
	"github.com/echocat/kubor/template/functions"
	"k8s.io/apimachinery/pkg/runtime"

	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	v1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

type Claim struct {
	GroupVersionKinds GroupVersionKinds `yaml:"gvks,omitempty" json:"gvks,omitempty"`
	SourceNamespaces  []string          `yaml:"namespaces,omitempty" json:"namespaces,omitempty"`

	// Values set using implicitly.
	Namespaces Namespaces `yaml:"-" json:"-"`
}

var (
	DefaultClaimedGroupVersionKinds = loadDefaultClaimedGroupVersionKinds()
)

func newClaim() Claim {
	return Claim{
		GroupVersionKinds: DefaultClaimedGroupVersionKinds,
		SourceNamespaces:  []string{"{{.GroupId}}"},
	}
}

func (instance Claim) evaluate(context interface{}) (Claim, error) {
	fail := func(n string, err error) (Claim, error) {
		return Claim{}, fmt.Errorf("cannot handle namespace '%s': %w", n, err)
	}
	result := instance
	result.Namespaces = make(Namespaces, len(result.SourceNamespaces))
	for i, source := range result.SourceNamespaces {
		if tmpl, err := functions.DefaultTemplateFactory().New(source, source); err != nil {
			return fail(source, err)
		} else if rendered, err := tmpl.ExecuteToString(context); err != nil {
			return fail(source, err)
		} else if err := result.Namespaces[i].Set(rendered); err != nil {
			return fail(source, err)
		}
	}
	return result, nil
}

func (instance Claim) Validate(reference ObjectReference) error {
	if !instance.Namespaces.Contains(reference.Namespace) {
		return fmt.Errorf("is in namespace %v; but claimed: %v", reference.Namespace, instance.Namespaces)
	}
	if !instance.GroupVersionKinds.Contains(reference.GroupVersionKind) {
		return fmt.Errorf("is group version kind %v; but claimed: %v", reference.GroupVersionKind, instance.GroupVersionKinds)
	}
	return nil
}

func loadDefaultClaimedGroupVersionKinds() GroupVersionKinds {
	s := runtime.NewScheme()
	s.AddKnownTypes(rbacv1.SchemeGroupVersion,
		&rbacv1.Role{},
		&rbacv1.RoleBinding{},
		&rbacv1.ClusterRole{},
		&rbacv1.ClusterRoleBinding{},
	)
	s.AddKnownTypes(v1.SchemeGroupVersion,
		&v1.Namespace{},
		&v1.Pod{},
		&v1.Service{},
		&v1.Secret{},
		&v1.ServiceAccount{},
		&v1.PersistentVolume{},
		&v1.PersistentVolumeClaim{},
		&v1.ConfigMap{},
	)
	s.AddKnownTypes(appsv1.SchemeGroupVersion,
		&appsv1.Deployment{},
		&appsv1.StatefulSet{},
		&appsv1.DaemonSet{},
	)
	s.AddKnownTypes(appsv1beta1.SchemeGroupVersion,
		&appsv1beta1.Deployment{},
		&appsv1beta1.StatefulSet{},
	)
	s.AddKnownTypes(appsv1beta2.SchemeGroupVersion,
		&appsv1beta2.Deployment{},
		&appsv1beta2.StatefulSet{},
		&appsv1beta2.DaemonSet{},
	)
	s.AddKnownTypes(apiextensions.SchemeGroupVersion,
		&apiextensions.CustomResourceDefinition{},
	)
	s.AddKnownTypes(apiextensionsv1.SchemeGroupVersion,
		&apiextensionsv1.CustomResourceDefinition{},
	)
	s.AddKnownTypes(apiextensionsv1beta1.SchemeGroupVersion,
		&apiextensionsv1beta1.CustomResourceDefinition{},
	)
	s.AddKnownTypes(extensionsv1beta1.SchemeGroupVersion,
		&extensionsv1beta1.Deployment{},
		&extensionsv1beta1.DaemonSet{},
		&extensionsv1beta1.Ingress{},
		&extensionsv1beta1.NetworkPolicy{},
	)
	s.AddKnownTypes(batchv1.SchemeGroupVersion,
		&batchv1.Job{},
	)
	s.AddKnownTypes(batchv1beta1.SchemeGroupVersion,
		&batchv1beta1.CronJob{},
	)
	s.AddKnownTypes(batchv2alpha1.SchemeGroupVersion,
		&batchv2alpha1.CronJob{},
	)
	return MapToGroupVersionKinds(s.AllKnownTypes())
}
