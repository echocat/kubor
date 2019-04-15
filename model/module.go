package model

type Module struct {
	GroupId    ResourceIdentifier `json:"groupId,omitempty" yaml:"groupId,omitempty"`
	ArtifactId ResourceIdentifier `json:"artifactId"        yaml:"artifactId"`
	Release    Release            `json:"release,omitempty" yaml:"release,omitempty"`
	Values     Values             `json:"values,omitempty"  yaml:"values,omitempty"`
}
