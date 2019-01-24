package command

import (
	"github.com/levertonai/kubor/common"
)

func init() {
	common.RegisterCliFactory(&ShowTemplate{})
}

type ShowTemplate struct {
	Command
}

func (instance *ShowTemplate) ConfigureCliCommands(context string, hc common.HasCommands) error {
	if context != "show" {
		return nil
	}
	cmd := hc.Command("template", "Show different values, commands etc. see sub-commands.")
	return common.ConfigureCliCommands("show/template", cmd)
}
