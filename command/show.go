package command

import (
	"github.com/echocat/kubor/common"
)

func init() {
	common.RegisterCliFactory(&Show{})
}

type Show struct {
	Command
}

func (instance *Show) ConfigureCliCommands(context string, hc common.HasCommands, version string) error {
	if context != "" {
		return nil
	}
	cmd := hc.Command("show", "Show different values, commands etc. see sub-commands.")
	return common.ConfigureCliCommands("show", cmd, version)
}
