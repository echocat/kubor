package template

import (
	. "github.com/onsi/gomega"
	"testing"
)

func Test_IsLiteral_detects_literal(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(IsLiteral("foo")).To(Equal(true))
	g.Expect(IsLiteral("foo{a}")).To(Equal(true))
	g.Expect(IsLiteral("foo{a}bar")).To(Equal(true))
	g.Expect(IsLiteral("foo{a}bar")).To(Equal(true))

	g.Expect(IsLiteral("{{foo}}")).To(Equal(false))
	g.Expect(IsLiteral("foo{{.}}bar")).To(Equal(false))
	g.Expect(IsLiteral("foo{{|}}bar")).To(Equal(false))
}
