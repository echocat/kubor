package command

import (
	"fmt"
	"github.com/urfave/cli"
	"io"
	"kubor/common"
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

	SourceHint bool
}

func (instance *Render) CreateCliCommands() ([]cli.Command, error) {
	return []cli.Command{{
		Name:   "render",
		Usage:  "Renders the instances of this project using the provided values.",
		Action: instance.ExecuteFromCli,
		Flags: []cli.Flag{
			cli.BoolTFlag{
				Name:        "sourceHint, sh",
				Usage:       "Prints to the output a comment which indicates where the rendered content organically comes from.",
				EnvVar:      "KUBOR_SOURCE_HINT",
				Destination: &instance.SourceHint,
			},
		},
	}}, nil
}

func (instance *Render) RunWithArguments(arguments CommandArguments) error {
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
