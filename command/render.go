package command

import (
	"fmt"
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

func (instance *Render) ConfigureCliCommands(hc common.HasCommands) error {
	cmd := hc.Command("render", "Renders the instances of this project using the provided values.").
		Action(instance.ExecuteFromCli)

	cmd.Flag("file", "If specified this file will be used as template and rendered with the current kubor environment"+
		" and printed to stdout. If empty all project template files will be used.").
		Short('f').
		PlaceHolder("<file>").
		Envar("KUBOR_TEMPLATE_FILE").
		Default(instance.TemplateFile).
		StringVar(&instance.TemplateFile)
	cmd.Flag("sourceHint", "Prints to the output a comment which indicates where the rendered content organically comes from.").
		Envar("KUBOR_SOURCE_HINT").
		Default(fmt.Sprint(instance.SourceHint)).
		BoolVar(&instance.SourceHint)

	return nil
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
