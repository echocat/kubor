package common

import (
	. "github.com/onsi/gomega"
	"testing"
)

func Test_GetObjectPathValue_works_with_nested_struct(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(GetObjectPathValue(objectA, "B", "C", "D", "E", "Value")).To(BeNil())
	g.Expect(GetObjectPathValue(objectA, "B", "C", "D", "Value")).To(Equal("helloD"))
	g.Expect(GetObjectPathValue(objectA, "B", "C", "Value")).To(Equal("helloC"))
	g.Expect(GetObjectPathValue(objectA, "B", "NotExisting")).To(BeNil())
	g.Expect(GetObjectPathValue(objectA, "B")).To(Equal(objectB))
	g.Expect(GetObjectPathValue(objectA, "B", "Value")).To(Equal("helloB"))
	g.Expect(GetObjectPathValue(objectA, "Value")).To(Equal("helloA"))
	g.Expect(GetObjectPathValue(objectA, "AnInt")).To(Equal(666))
}

func Test_GetObjectPathValue_works_with_nested_maps(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(GetObjectPathValue(mapA, "B", "C", "D", "E", "Value")).To(BeNil())
	g.Expect(GetObjectPathValue(mapA, "B", "C", "D", "Value")).To(Equal("helloD"))
	g.Expect(GetObjectPathValue(mapA, "B", "C", "Value")).To(Equal("helloC"))
	g.Expect(GetObjectPathValue(mapA, "B", "NotExisting")).To(BeNil())
	g.Expect(GetObjectPathValue(mapA, "B", "Value")).To(Equal("helloB"))
	g.Expect(GetObjectPathValue(mapA, "Value")).To(Equal("helloA"))
	g.Expect(GetObjectPathValue(mapA, "AnInt")).To(Equal(666))
}
