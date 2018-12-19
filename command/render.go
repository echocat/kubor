package command

import (
	"fmt"
	"github.com/urfave/cli"
	"io"
	"kubor/common"
	"os"
	"strings"
)

const (
	sourceHintTemplate = `######################################################################
# %s
######################################################################
`
)

func init() {
	cmd := &Render{}
	cmd.Parent = cmd
	RegisterInitializable(cmd)
	common.RegisterCliFactory(cmd)
}

type Render struct {
	Command

	TemplateFile string
	SourceHint   bool
}

func (instance *Render) CreateCliCommands(context string) ([]cli.Command, error) {
	if context != "" {
		return nil, nil
	}
	return []cli.Command{{
		Name:   "render",
		Usage:  "Renders the instances of this project using the provided values.",
		Action: instance.ExecuteFromCli,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "file, f",
				Usage: "If specified this file will be used as template and rendered with the current kubor environment\n" +
					"\tand printed to stdout. If empty all project template files will be used.",
				EnvVar:      "KUBOR_TEMPLATE_FILE",
				Destination: &instance.TemplateFile,
			},
			cli.BoolTFlag{
				Name: "sourceHint",
				Usage: "Prints to the output a comment which indicates where the rendered content organically\n" +
					"\tcomes from. This will be ignored if --file is used.",
				EnvVar:      "KUBOR_SOURCE_HINT",
				Destination: &instance.SourceHint,
			},
		},
	}}, nil
}

func (instance *Render) RunWithArguments(arguments CommandArguments) error {
	if instance.TemplateFile != "" {
		return arguments.Project.RenderedTemplateFile(instance.TemplateFile, os.Stdout)
	}

	cp, err := arguments.Project.RenderedTemplatesProvider()
	if err != nil {
		return err
	}

	var name string
	var content []byte
	first := true

	for name, content, err = cp(); err == nil; name, content, err = cp() {
		trimmed := strings.TrimSpace(string(content))
		if len(trimmed) > 0 {
			if first {
				first = false
			} else {
				fmt.Print("\n---\n")
			}
			if instance.SourceHint {
				fmt.Printf(sourceHintTemplate, name)
			}
			fmt.Print(trimmed)
		}
	}
	if err == io.EOF {
		return nil
	}

	return err
}
