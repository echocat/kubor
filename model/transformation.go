package model

type Transformation struct {
	Enabled  *bool   `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Argument *string `yaml:"argument,omitempty" json:"argument,omitempty"`
}

func (instance Transformation) Merge(with Transformation) Transformation {
	result := instance

	if v := with.Enabled; v != nil {
		result.Enabled = v
	}
	if v := with.Argument; v != nil {
		result.Argument = v
	}

	return result
}

func (instance Transformation) IsEnabled(def bool) bool {
	if v := instance.Enabled; v != nil {
		return *v
	}
	return def
}

func (instance Transformation) ArgumentAsString() string {
	if v := instance.Argument; v != nil {
		return *v
	}
	return "nil"
}
