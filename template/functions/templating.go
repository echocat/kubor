package functions

import (
	"fmt"
	"github.com/levertonai/kubor/template"
)

var FuncRender = Function{
	Description: "Renders the <template> using <data> as regular Golang template.",
	Parameters: Parameters{{
		Name:        "data",
		Description: "The data that could be accessed while the rendering the content of the provided <template>.",
	}, {
		Name:        "template",
		Description: "The actual template which should be rendered using the provided <data>.",
	}},
	Returns: Return{
		Description: "The rendered content.",
	},
}.MustWithFunc(func(context template.ExecutionContext, data interface{}, template string) (string, error) {
	contextTmpl := context.GetTemplate()
	if tmpl, err := context.GetFactory().New(contextTmpl.GetSourceName(), template); err != nil {
		return "", fmt.Errorf("%s: cannot create parse template: %v", contextTmpl.GetSource(), err)
	} else if result, err := tmpl.ExecuteToString(data); err != nil {
		return "", fmt.Errorf("%s: cannot evaluate parse template: %v", tmpl.GetSource(), err)
	} else {
		return result, nil
	}
})

var FuncInclude = Function{
	Description: "Takes the given <file> and renders the contained template using <data> as regular Golang template.",
	Parameters: Parameters{{
		Name:        "file",
		Description: "The actual template file which should be rendered using the provided <data>.",
	}, {
		Name:        "data",
		Description: "The data that could be accessed while the rendering the content of the template of the provided <file>.",
	}},
	Returns: Return{
		Description: "The rendered content.",
	},
}.MustWithFunc(func(context template.ExecutionContext, file string, data interface{}) (string, error) {
	if resolved, err := resolvePathOfContext(context, file); err != nil {
		return "", err
	} else if tmpl, err := context.GetFactory().NewFromFile(resolved); err != nil {
		return "", fmt.Errorf("%s: cannot create parse template: %v", resolved, err)
	} else if result, err := tmpl.ExecuteToString(data); err != nil {
		return "", fmt.Errorf("%s: cannot evaluate parse template: %v", tmpl.GetSource(), err)
	} else {
		return result, nil
	}
})

var FuncSourceFile = Function{
	Returns: Return{
		Description: "The filename which is the source of this rendered template if any.",
	},
}.MustWithFunc(func(context template.ExecutionContext) *string {
	return context.GetTemplate().GetSourceFile()
})

var FuncSourceName = Function{
	Returns: Return{
		Description: "The name this rendered template. Could be a short name or a filename.",
	},
}.MustWithFunc(func(context template.ExecutionContext) string {
	return context.GetTemplate().GetSourceName()
})

var FuncsTemplating = Functions{
	"render":     FuncRender,
	"include":    FuncInclude,
	"sourceFile": FuncSourceFile,
	"sourceName": FuncSourceName,
}
var CategoryTemplating = Category{
	Functions: FuncsTemplating,
}
