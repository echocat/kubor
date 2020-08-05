package support

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NormalizeLabelValue(t *testing.T) {
	cases := []struct {
		given    string
		expected string
	}{
		{"_-foo-bar_-", "foo-bar"},
		{"x_-foo-bar_-", "x_-foo-bar"},
		{"_-foo-bar_-x", "foo-bar_-x"},
		{"foo_bar", "foo_bar"},
		{"foo.bar", "foo.bar"},
		{"foo:bar", "foo_bar"},
		{"foo/bar", "foo_bar"},
		{"foo,bar", "foo_bar"},
		{"abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyzx123456789", "abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyzx"},
	}
	for _, c := range cases {
		t.Run(c.given, func(t *testing.T) {
			assert.Equal(t, NormalizeLabelValue(c.given), c.expected)
		})
	}
}
