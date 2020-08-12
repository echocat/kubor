package model

type Transformations map[TransformationName]Transformation

func NewTransformations() Transformations {
	return Transformations{}
}

func (instance Transformations) Get(name TransformationName) (result Transformation, err error) {
	if v := instance; v != nil {
		result = v[name]
	}
	return
}
