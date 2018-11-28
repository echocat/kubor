package common

import (
	. "github.com/onsi/gomega"
	"testing"
)

func Test_MustEvaluateExpression_works_with_nested_struct(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(MustEvaluateExpression("B.C.D.E.Value", objectA)).To(BeNil())
	g.Expect(MustEvaluateExpression("B.C.D.Value", objectA)).To(Equal("helloD"))
	g.Expect(MustEvaluateExpression("B.C.Value", objectA)).To(Equal("helloC"))
	g.Expect(MustEvaluateExpression("B.NotExisting", objectA)).To(BeNil())
	g.Expect(MustEvaluateExpression("B", objectA)).To(Equal(objectB))
	g.Expect(MustEvaluateExpression("B.Value", objectA)).To(Equal("helloB"))
	g.Expect(MustEvaluateExpression("Value", objectA)).To(Equal("helloA"))
	g.Expect(MustEvaluateExpression("AnInt", objectA)).To(Equal(666))
}

func Test_MustEvaluateExpression_works_with_nested_maps(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(MustEvaluateExpression("B.C.D.E.Value", mapA)).To(BeNil())
	g.Expect(MustEvaluateExpression("B.C.D.Value", mapA)).To(Equal("helloD"))
	g.Expect(MustEvaluateExpression("B.C.Value", mapA)).To(Equal("helloC"))
	g.Expect(MustEvaluateExpression("B.NotExisting", mapA)).To(BeNil())
	g.Expect(MustEvaluateExpression("B.Value", mapA)).To(Equal("helloB"))
	g.Expect(MustEvaluateExpression("Value", mapA)).To(Equal("helloA"))
	g.Expect(MustEvaluateExpression("AnInt", mapA)).To(Equal(666))
}

func Test_MustEvaluateExpression_panics_on_invalid_expression(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(func() {
		MustEvaluateExpression("..", mapA)
	}).To(Panic())
}
