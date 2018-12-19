package command

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"kubor/common"
	"kubor/template/functions"
	"os"
	"regexp"
)

const (
	showTemplateFunctionsTemplate = `List of available function within kubor rending context by category:

{{range $category, $functions := .FunctionsByCategories -}}
{{ $category | upper }}:{{range $i, $function := $functions}}
  {{- "\n  {{ " -}}{{- $function.GetName -}}
    {{- range $i, $parameter := $function.Parameters }}
       {{- if $parameter.VarArg }} [<{{ $parameter.GetName }}> ...]
       {{- else }} <{{ $parameter.GetName }}>{{ end -}}
    {{- end -}}{{- " }}" }}
    {{- if $function.Description }}
    {{ $function.Description | replace "\n" " " | warpCustom 120 "\n    " false }}{{ end -}} 
    {{- range $i, $parameter := $function.Parameters }}
    - {{ $parameter.GetName }}: {{ $parameter.GetType -}}
      {{- if $parameter.Description }} - {{ $parameter.Description | replace "\n" " " | warpCustom 120 "\n      " false  }}{{ end -}}
    {{- end }}
    Returns: {{ $function.Returns.GetType -}}
      {{- if $function.Returns.Description }} - {{ $function.Returns.Description | replace "\n" " " | warpCustom 120 "\n      " false  }}{{ end }}
{{ end }}
{{ end -}}
NOTES:
  Please refer https://golang.org/pkg/text/template/ for more information about the template language of kubor.
`
)

type ShowTemplateFunctionsOutput string

func (instance *ShowTemplateFunctionsOutput) Set(plain string) error {
	if plain != "json" && plain != "yaml" && plain != "text" {
		return fmt.Errorf("unsupported output format: %s", plain)
	}
	*instance = ShowTemplateFunctionsOutput(plain)
	return nil
}

func (instance ShowTemplateFunctionsOutput) String() string {
	return string(instance)
}

func init() {
	cmd := &ShowTemplateFunctions{
		Output: ShowTemplateFunctionsOutput("text"),
	}
	common.RegisterCliFactory(cmd)
}

type ShowTemplateFunctions struct {
	Output ShowTemplateFunctionsOutput
}

func (instance *ShowTemplateFunctions) CreateCliCommands(context string) ([]cli.Command, error) {
	if context != "show" {
		return nil, nil
	}
	return []cli.Command{{
		Name:      "templateFunctions",
		ArgsUsage: "[function name regexp]",
		Usage:     "Shows a list of all available template functions.",
		UsageText: "If [function name regexp] specified only functions that matches at least one of the given patterns.",
		Action:    instance.ExecuteFromCli,
		Flags: []cli.Flag{
			cli.GenericFlag{
				Name:  "output, o",
				Usage: "Specifies how to render the output.",
				Value: &instance.Output,
			},
		},
	}}, nil
}

func (instance *ShowTemplateFunctions) createPredicate(c *cli.Context) (func(functions.Function) bool, error) {
	regexps := make([]*regexp.Regexp, c.NArg())
	for i, pattern := range c.Args() {
		if r, err := regexp.Compile(pattern); err != nil {
			return nil, fmt.Errorf("problems while compile function regexp pattern: %v", err)
		} else {
			regexps[i] = r
		}
	}
	if len(regexps) <= 0 {
		return func(functions.Function) bool {
			return true
		}, nil
	}
	return func(in functions.Function) bool {
		for _, r := range regexps {
			if r.MatchString(in.GetName()) {
				return true
			}
		}
		return false
	}, nil
}

func (instance *ShowTemplateFunctions) ExecuteFromCli(c *cli.Context) error {
	predicate, err := instance.createPredicate(c)
	if err != nil {
		return err
	}
	all := functions.GlobalRegistry().GetAll()

	context := instance.newShowTemplateFunctionsContext(all, predicate)

	switch instance.Output {
	case ShowTemplateFunctionsOutput("yaml"):
		return context.renderYaml()
	case ShowTemplateFunctionsOutput("json"):
		return context.renderJson()
	default:
		return context.renderText()
	}
}

func (instance *ShowTemplateFunctions) newShowTemplateFunctionsContext(in []functions.Function, predicate func(functions.Function) bool) showTemplateFunctionsContext {
	functionsByCategories := map[string][]functions.Function{}
	var lengthOfLongestName int
	var lengthOfLongestCategory int
	for _, function := range in {
		if predicate(function) {
			name := function.GetName()
			category := function.GetCategory()
			if lengthOfLongestName < len(name) {
				lengthOfLongestName = len(name)
			}
			if lengthOfLongestCategory < len(category) {
				lengthOfLongestCategory = len(category)
			}

			if existing, ok := functionsByCategories[category]; ok {
				functionsByCategories[category] = append(existing, function)
			} else {
				functionsByCategories[category] = []functions.Function{function}
			}
		}
	}
	return showTemplateFunctionsContext{
		FunctionsByCategories:   functionsByCategories,
		LengthOfLongestName:     lengthOfLongestName,
		LengthOfLongestCategory: lengthOfLongestCategory,
	}
}

type showTemplateFunctionsContext struct {
	FunctionsByCategories   map[string][]functions.Function
	LengthOfLongestName     int
	LengthOfLongestCategory int
}

func (instance showTemplateFunctionsContext) renderYaml() error {
	enc := yaml.NewEncoder(os.Stdout)
	return enc.Encode(instance.FunctionsByCategories)
}

func (instance showTemplateFunctionsContext) renderJson() error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(instance.FunctionsByCategories)
}

func (instance showTemplateFunctionsContext) renderText() error {
	if tmpl, err := functions.GlobalTemplateFactory().New(
		"show_templateFunctions.go",
		showTemplateFunctionsTemplate,
	); err != nil {
		return err
	} else {
		return tmpl.Execute(instance, os.Stdout)
	}
}
