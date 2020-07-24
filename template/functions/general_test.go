package functions

import (
	"gotest.tools/assert"
	"testing"
)

func Test_FuncMap(t *testing.T) {
	assert.Equal(t, mustExecuteTemplate(t, `{{ map }}`, nil), "map[]")
	assert.Equal(t, mustExecuteTemplate(t, `{{ (map "foo" 666).foo }}`, nil), "666")
	assert.Equal(t, mustExecuteTemplate(t, `{{ $m := map "foo" 666 "bar" 123 "xyz" . }}
	                                        {{- $m.foo}},{{$m.bar}},{{$m.xyz}}`, "hello"), "666,123,hello")

	_, err := executeTemplate(t, `{{map "foo"}}`, nil)
	assert.ErrorContains(t, err, "expect always a key to value pair, this means the amount of parameters needs to be dividable by two, but got: 1")
}

func Test_FuncSlice(t *testing.T) {
	assert.Equal(t, mustExecuteTemplate(t, `{{ slice }}`, nil), "[]")
	assert.Equal(t, mustExecuteTemplate(t, `{{ slice "foo" 666 }}`, nil), "[foo 666]")
}

func Test_FuncContains(t *testing.T) {
	str := "abc"
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "b" }}`, str), "true")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "x" }}`, str), "false")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "b" }}`, &str), "true")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "x" }}`, &str), "false")

	m := map[interface{}]interface{}{
		1:   2,
		"a": "b",
	}
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains 1 }}`, m), "true")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains 2 }}`, m), "false")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "a" }}`, &m), "true")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "x" }}`, &m), "false")

	s := struct {
		a string
		b string
	}{
		a: "foo",
		b: "bar",
	}
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "a" }}`, s), "true")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "x" }}`, s), "false")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "a" }}`, &s), "true")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "x" }}`, &s), "false")

	sl := []interface{}{"a", "b", 1, 2}
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "a" }}`, sl), "true")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains 1 }}`, sl), "true")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "x" }}`, sl), "false")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains 666 }}`, sl), "false")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "a" }}`, &sl), "true")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains 1 }}`, &sl), "true")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "x" }}`, &sl), "false")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains 666 }}`, &sl), "false")

	a := [4]interface{}{"a", "b", 1, 2}
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "a" }}`, a), "true")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains 1 }}`, a), "true")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "x" }}`, a), "false")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains 666 }}`, a), "false")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "a" }}`, &a), "true")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains 1 }}`, &a), "true")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains "x" }}`, &a), "false")
	assert.Equal(t, mustExecuteTemplate(t, `{{ . | contains 666 }}`, &a), "false")
}
