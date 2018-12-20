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

{{range $categoryName, $category := .Categories -}}
{{ $categoryName | upper }}:{{range $functionName, $function := $category.Functions}}
  {{- "\n  {{ " -}}{{- $functionName -}}
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

type ShowTemplateFulltextTerm struct {
	Regexp *regexp.Regexp
}

func (instance *ShowTemplateFulltextTerm) Set(plain string) error {
	if plain == "" {
		(*instance).Regexp = nil
		return nil
	} else if r, err := regexp.Compile("(?i)" + plain); err != nil {
		return err
	} else {
		(*instance).Regexp = r
		return nil
	}
}

func (instance ShowTemplateFulltextTerm) String() string {
	r := instance.Regexp
	if r == nil {
		return ""
	}
	return r.String()
}

func init() {
	cmd := &ShowTemplateFunctions{
		Output: ShowTemplateFunctionsOutput("text"),
	}
	common.RegisterCliFactory(cmd)
}

type ShowTemplateFunctions struct {
	Output             ShowTemplateFunctionsOutput
	FulltextSearchTerm ShowTemplateFulltextTerm
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
			}, cli.GenericFlag{
				Name:  "fulltextSearch, fulltext, s",
				Usage: "Specifies a term that needs to be matched in any property of a function (name, description, ...)",
				Value: &instance.FulltextSearchTerm,
			},
		},
	}}, nil
}

func (instance *ShowTemplateFunctions) createPredicate(c *cli.Context) (func(name string) bool, error) {
	regexps := make([]*regexp.Regexp, c.NArg())
	for i, pattern := range c.Args() {
		if r, err := regexp.Compile(pattern); err != nil {
			return nil, fmt.Errorf("problems while compile function regexp pattern: %v", err)
		} else {
			regexps[i] = r
		}
	}
	if len(regexps) <= 0 {
		return func(name string) bool {
			return true
		}, nil
	}
	return func(name string) bool {
		for _, r := range regexps {
			if r.MatchString(name) {
				return true
			}
		}
		return false
	}, nil
}

func (instance *ShowTemplateFunctions) createFulltextPredicate(c *cli.Context) func(functionName string, function functions.Function) bool {
	term := instance.FulltextSearchTerm.Regexp
	if term == nil {
		return func(functionName string, function functions.Function) bool {
			return true
		}
	}
	return func(functionName string, function functions.Function) bool {
		if term.FindStringIndex(functionName) != nil {
			return true
		}
		if function.MatchesFulltextSearch(term) {
			return true
		}
		return false
	}
}

func (instance *ShowTemplateFunctions) ExecuteFromCli(c *cli.Context) error {
	namePredicate, err := instance.createPredicate(c)
	if err != nil {
		return err
	}
	fulltextPredicate := instance.createFulltextPredicate(c)

	context := instance.newShowTemplateFunctionsContext(namePredicate, fulltextPredicate, functions.CategoriesDefault)

	switch instance.Output {
	case ShowTemplateFunctionsOutput("yaml"):
		return context.renderYaml()
	case ShowTemplateFunctionsOutput("json"):
		return context.renderJson()
	default:
		return context.renderText()
	}
}

func (instance *ShowTemplateFunctions) newShowTemplateFunctionsContext(
	namePredicate func(name string) bool,
	fulltextPredicate func(functionName string, function functions.Function) bool,
	categories functions.Categories,
) showTemplateFunctionsContext {
	targetCategories := functions.Categories{}
	for categoryName, category := range categories {
		for functionName, function := range category.Functions {
			if namePredicate(functionName) && fulltextPredicate(functionName, function) {
				targetCategories = targetCategories.WithFunction(categoryName, functionName, function)
			}
		}
	}
	return showTemplateFunctionsContext{
		Categories: targetCategories,
	}
}

type showTemplateFunctionsContext struct {
	Categories functions.Categories
}

func (instance showTemplateFunctionsContext) renderYaml() error {
	enc := yaml.NewEncoder(os.Stdout)
	return enc.Encode(instance.Categories)
}

func (instance showTemplateFunctionsContext) renderJson() error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(instance.Categories)
}

func (instance showTemplateFunctionsContext) renderText() error {
	if tmpl, err := functions.DefaultTemplateFactory().New(
		"show_templateFunctions.go",
		showTemplateFunctionsTemplate,
	); err != nil {
		return err
	} else {
		return tmpl.Execute(instance, os.Stdout)
	}
}
