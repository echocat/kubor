package transformation

import (
	"github.com/echocat/kubor/kubernetes/support"
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func init() {
	Default.MustRegisterUpdateFunc("apply-labels", func(project *model.Project, _ unstructured.Unstructured, target *unstructured.Unstructured, argument string) error {
		return ensureKuborLabels(project, target, argument)
	})
	Default.MustRegisterCreateFunc("apply-labels", ensureKuborLabels)
}

func ensureKuborLabels(project *model.Project, target *unstructured.Unstructured, _ string) error {
	if err := ensureKuborLabelsOfPath(project, target, "metadata", "labels"); err != nil {
		return err
	}
	if err := ensureKuborLabelsOfPath(project, target, "spec", "template", "metadata", "labels"); err != nil {
		return err
	}
	return nil
}

func ensureKuborLabelsOfPath(project *model.Project, target *unstructured.Unstructured, fields ...string) error {
	pl := project.Labels
	labels, _, err := unstructured.NestedStringMap(target.Object, fields...)
	if err != nil {
		return err
	}
	if labels == nil {
		labels = make(map[string]string)
	}

	ensureKuborLabel(&labels, pl.GroupId, project.GroupId.String())
	ensureKuborLabel(&labels, pl.ArtifactId, project.ArtifactId.String())
	ensureKuborLabel(&labels, pl.Release, support.NormalizeLabelValue(project.Release))

	return unstructured.SetNestedStringMap(target.Object, labels, fields...)
}

func ensureKuborLabel(labels *map[string]string, label model.Label, value string) {
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
