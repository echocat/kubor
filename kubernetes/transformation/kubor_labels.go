package transformation

import (
	"github.com/echocat/kubor/kubernetes/support"
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func init() {
	RegisterUpdateTransformationFunc(func(project *model.Project, _ unstructured.Unstructured, target *unstructured.Unstructured) error {
		return ensureKuborLabels(project, target)
	})
	RegisterCreateTransformationFunc(ensureKuborLabels)
}

func ensureKuborLabels(project *model.Project, target *unstructured.Unstructured) error {
	pl := project.Labels
	labels := target.GetLabels()

	ensureKuborLabel(&labels, pl.GroupId, project.GroupId)
	ensureKuborLabel(&labels, pl.ArtifactId, project.ArtifactId)
	ensureKuborLabel(&labels, pl.Release, project.Release)

	target.SetLabels(labels)

	return nil
}

func ensureKuborLabel(labels *map[string]string, label model.Label, value string) {
	switch label.Action {
	case model.LabelActionDrop:
		delete(*labels, label.Name.String())
	case model.LabelActionSet:
		(*labels)[label.Name.String()] = support.NormalizeLabelValue(value)
	case model.LabelActionSetIfAbsent:
		if _, exist := (*labels)[label.Name.String()]; !exist {
			(*labels)[label.Name.String()] = support.NormalizeLabelValue(value)
		}
	case model.LabelActionSetIfExists:
		if _, exist := (*labels)[label.Name.String()]; exist {
			(*labels)[label.Name.String()] = support.NormalizeLabelValue(value)
		}
	}
}
