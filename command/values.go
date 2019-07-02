package command

import (
	"github.com/alecthomas/kingpin"
	"github.com/echocat/kubor/common"
	"gopkg.in/yaml.v2"
	"os"
)

func init() {
	cmd := &Values{}
	cmd.Parent = cmd
	RegisterInitializable(cmd)
	common.RegisterCliFactory(cmd)
}

type Values struct {
	Command
}

func (instance *Values) ConfigureCliCommands(context string, hc common.HasCommands, version string) error {
	if context != "" {
		return nil
	}
	hc.Command("values", "Get the aggregated values used by the project based on the given parameters.").
		Action(func(context *kingpin.ParseContext) error {
			return instance.Run()
		})
	return nil
}

func (instance *Values) RunWithArguments(arguments Arguments) error {
	enc := yaml.NewEncoder(os.Stdout)
	return enc.Encode(arguments.Project.Values)
}
