package transformation

import (
	"github.com/echocat/kubor/kubernetes/support"
	"github.com/echocat/kubor/model"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const kuborLabelName = model.TransformationName("apply-labels")

func init() {
	t := ensureKuborLabels{}
	Default.MustRegisterUpdate(&t)
	Default.MustRegisterUpdate(&t)
}

var NamespaceGvks = model.BuildGroupVersionKinds(v1.SchemeGroupVersion, &v1.Namespace{}).Build()

type ensureKuborLabels struct{}

func (instance *ensureKuborLabels) GetName() model.TransformationName {
	return kuborLabelName
}

func (instance *ensureKuborLabels) GetPriority() int32 {
	return 1_000_000_001
}

func (instance *ensureKuborLabels) DefaultEnabled(target *unstructured.Unstructured) bool {
	return !NamespaceGvks.Contains(model.GroupVersionKind(target.GroupVersionKind()))
}

func (instance *ensureKuborLabels) TransformForUpdate(p *model.Project, _ unstructured.Unstructured, target *unstructured.Unstructured, argument *string) error {
	return instance.TransformForCreate(p, target, argument)
}

func (instance *ensureKuborLabels) TransformForCreate(project *model.Project, target *unstructured.Unstructured, _ *string) error {
	if err := instance.ensureKuborLabelsOfPath(project, target, "metadata", "labels"); err != nil {
		return err
	}
	if err := instance.ensureKuborLabelsOfPath(project, target, "spec", "template", "metadata", "labels"); err != nil {
		return err
	}
	return nil
}

func (instance ensureKuborLabels) ensureKuborLabelsOfPath(project *model.Project, target *unstructured.Unstructured, fields ...string) error {
	pl := project.Labels
	labels, _, err := unstructured.NestedStringMap(target.Object, fields...)
	if err != nil {
		return err
	}
	if labels == nil {
		labels = make(map[string]string)
	}

	instance.ensureKuborLabel(&labels, pl.GroupId, project.GroupId.String())
	instance.ensureKuborLabel(&labels, pl.ArtifactId, project.ArtifactId.String())
	instance.ensureKuborLabel(&labels, pl.Release, support.NormalizeLabelValue(project.Release))

	return unstructured.SetNestedStringMap(target.Object, labels, fields...)
}

func (instance ensureKuborLabels) ensureKuborLabel(labels *map[string]string, label model.Label, value string) {
	switch label.Action {
	case model.LabelActionDrop:
		delete(*labels, label.Name.String())
	case model.LabelActionSet:
		(*labels)[label.Name.String()] = value
	case model.LabelActionSetIfAbsent:
		if _, exist := (*labels)[label.Name.String()]; !exist {
			(*labels)[label.Name.String()] = value
		}
	case model.LabelActionSetIfExists:
		if _, exist := (*labels)[label.Name.String()]; exist {
			(*labels)[label.Name.String()] = value
		}
	}
}
