package command

import (
	"github.com/echocat/kubor/common"
	"github.com/echocat/kubor/kubernetes"
)

func init() {
	cmd := &Delete{}
	cmd.Parent = cmd
	RegisterInitializable(cmd)
	common.RegisterCliFactory(cmd)
}

type Delete struct {
	Command
}

func (instance *Delete) ConfigureCliCommands(context string, hc common.HasCommands, _ string) error {
	if context != "" {
		return nil
	}

	hc.Command("delete", "Will delete all resources which matches the current"+
		" project's groupId and artifactId in the configured claim.").
		Action(instance.ExecuteFromCli)
	return nil
}

func (instance *Delete) RunWithArguments(arguments Arguments) error {
	ct, err := kubernetes.NewCleanupTask(arguments.Project, arguments.DynamicClient, kubernetes.CleanupModeDelete)
	if err != nil {
		return err
	}
	return ct.Execute()
}
