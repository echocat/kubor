package command

import (
	"github.com/alecthomas/kingpin"
	"github.com/levertonai/kubor/common"
	"gopkg.in/yaml.v2"
	"os"
)

func init() {
	cmd := &Claims{}
	cmd.Parent = cmd
	RegisterInitializable(cmd)
	common.RegisterCliFactory(cmd)
}

type Claims struct {
	Command
}

func (instance *Claims) ConfigureCliCommands(context string, hc common.HasCommands) error {
	if context != "" {
		return nil
	}
	hc.Command("claims", "Get the aggregated claims used by the project based on the given parameters.").
		Action(func(context *kingpin.ParseContext) error {
			return instance.Run()
		})
	return nil
}

func (instance *Claims) RunWithArguments(arguments Arguments) error {
	enc := yaml.NewEncoder(os.Stdout)
	return enc.Encode(arguments.Project.Claims)
}
