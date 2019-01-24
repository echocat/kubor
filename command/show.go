package command

import (
	"github.com/levertonai/kubor/common"
)

func init() {
	common.RegisterCliFactory(&Show{})
}

type Show struct {
	Command
}

func (instance *Show) ConfigureCliCommands(context string, hc common.HasCommands) error {
	if context != "" {
		return nil
	}
	cmd := hc.Command("show", "Show different values, commands etc. see sub-commands.")
	return common.ConfigureCliCommands("show", cmd)
}
