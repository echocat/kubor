package model

import "github.com/levertonai/kubor/template/functions"

var (
	templateFactory = functions.DefaultTemplateFactory()
)

func evaluateTemplate(name, source string, with interface{}) (string, error) {
	if tmpl, err := templateFactory.New(name, source); err != nil {
		return "", err
	} else if result, err := tmpl.ExecuteToString(with); err != nil {
		return "", err
	} else {
		return result, nil
	}
}
