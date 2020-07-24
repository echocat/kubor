package functions

import (
	"github.com/echocat/kubor/template"
	"testing"
)

func templateBy(t *testing.T, source string) template.Template {
	result, err := DefaultTemplateFactory().New(t.Name(), source)
	if err != nil {
		t.Errorf("cannot parse given source: %v", err)
	}
	return result
}

func executeTemplate(t *testing.T, source string, context interface{}) (string, error) {
	tmpl := templateBy(t, source)
	return tmpl.ExecuteToString(context)
}

func mustExecuteTemplate(t *testing.T, source string, context interface{}) string {
	result, err := executeTemplate(t, source, context)
	if err != nil {
		t.Errorf("cannot execute template source: %v", err)
	}
	return result
}
