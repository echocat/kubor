package functions

import (
	"github.com/echocat/kubor/kubernetes/support"
)

var FuncNormalizeLabelValue = Function{
	Description: "Takes the given string and normalize it to fit into a kubernetes label value.",
	Parameters: Parameters{{
		Name: "source",
	}},
	Returns: Return{
		Description: "Normalized string of <source>.",
	},
}.MustWithFunc(func(source string) string {
	return support.NormalizeLabelValue(source)
})

var FuncsKubernetes = Functions{
	"normalizeLabelValue": FuncNormalizeLabelValue,
}
var CategoryKubernetes = Category{
	Functions: FuncsKubernetes,
}
