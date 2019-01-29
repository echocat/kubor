package format

type Value interface {
	Formatted() *string
	Interface() interface{}
	LessThan(Value) bool
}

type StringValue struct {
	Content *string
}

func (instance StringValue) Formatted() *string {
	return instance.Content
}

func (instance StringValue) Interface() interface{} {
	return instance.Content
}

func (instance StringValue) LessThan(what Value) bool {
	if sv, ok := what.(*StringValue); ok {
		what = *sv
	}
	if sv, ok := what.(StringValue); !ok {
		return false
	} else if instance.Content == nil && sv.Content == nil {
		return false
	} else if instance.Content == nil {
		return true
	} else if sv.Content == nil {
		return false
	} else {
		return *instance.Content < *sv.Content
	}
}

type Int32Value struct {
	Content *int32
}

func (instance Int32Value) Formatted() *string {
	return instance.Content
}

func (instance Int32Value) Interface() interface{} {
	return instance.Content
}

func (instance Int32Value) LessThan(what Value) bool {
	if sv, ok := what.(*Int32Value); ok {
		what = *sv
	}
	if sv, ok := what.(Int32Value); !ok {
		return false
	} else if instance.Content == nil && sv.Content == nil {
		return false
	} else if instance.Content == nil {
		return true
	} else if sv.Content == nil {
		return false
	} else {
		return *instance.Content < *sv.Content
	}
}
