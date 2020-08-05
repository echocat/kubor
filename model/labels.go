package model

const (
	LabelGroupId    = "kubor.echocat.org/group-id"
	LabelArtifactId = "kubor.echocat.org/artifact-id"
	LabelRelease    = "kubor.echocat.org/release"
)

type Labels struct {
	GroupId    Label `yaml:"groupId,omitempty" json:"groupId,omitempty"`
	ArtifactId Label `yaml:"artifactId,omitempty" json:"artifactId,omitempty"`
	Release    Label `yaml:"release,omitempty" json:"release,omitempty"`
}

func newLabels() Labels {
	return Labels{
		GroupId:    Label{LabelGroupId, LabelActionSetIfAbsent},
		ArtifactId: Label{LabelArtifactId, LabelActionSetIfAbsent},
		Release:    Label{LabelRelease, LabelActionSetIfAbsent},
	}
}
