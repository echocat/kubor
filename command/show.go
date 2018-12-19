package command

import (
	"github.com/urfave/cli"
	"kubor/common"
)

func init() {
	common.RegisterCliFactory(&Show{})
}

type Show struct {
	Command
}

func (instance *Show) CreateCliCommands(context string) ([]cli.Command, error) {
	if context != "" {
		return nil, nil
	}
	if commands, err := common.CreateCliCommands("show"); err != nil {
		return nil, err
	} else {
		return []cli.Command{{
			Name:        "show",
			Usage:       "Show different values, commands etc. see sub-commands. ",
			Subcommands: commands,
		}}, nil
	}
}
