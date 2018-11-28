package command

import (
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"kubor/common"
	"kubor/model"
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

func (instance *Values) CreateCliCommands() ([]cli.Command, error) {
	return []cli.Command{{
		Name:   "values",
		Usage:  "Get the current project using the provided properties",
		Action: instance.ExecuteFromCli,
	}}, nil
}

func (instance *Values) RunForProject(project model.Project) error {
	enc := yaml.NewEncoder(os.Stdout)
	return enc.Encode(project.Values)
}
