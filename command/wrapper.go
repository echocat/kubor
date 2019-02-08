package command

import (
	"github.com/alecthomas/kingpin"
	"github.com/levertonai/kubor/common"
	"github.com/levertonai/kubor/wrapper"
)

func init() {
	common.RegisterCliFactory(&Wrapper{})
}

type Wrapper struct {
	Command

	TargetDirectory string
	Version         string
}

func (instance *Wrapper) ConfigureCliCommands(context string, hc common.HasCommands, version string) error {
	if context != "" {
		return nil
	}
	wrapperCmd := hc.Command("wrapper", "Wrapper related commands.")

	wrapperCmd.Flag("targetDirectory", "Directory where to install the wrapper to.").
		Default(".").
		StringVar(&instance.TargetDirectory)
	wrapperCmd.Flag("version", "Version of the wrapper to be use.").
		Default(version).
		StringVar(&instance.Version)

	wrapperCmd.Command("install", "Will install the wrapper - will fail if there is already a wrapper.").
		Action(instance.executeInstall)
	wrapperCmd.Command("update", "Will update an existing wrapper - will fail if there is not already a wrapper.").
		Action(instance.executeUpdate)
	wrapperCmd.Command("ensure", "Will ensure that a wrapper exist.").
		Action(instance.executeEnsure)

	return nil
}

func (instance *Wrapper) executeInstall(_ *kingpin.ParseContext) error {
	return wrapper.Write(instance.TargetDirectory, instance.Version, wrapper.WoCreateOnly)
}

func (instance *Wrapper) executeUpdate(_ *kingpin.ParseContext) error {
	return wrapper.Write(instance.TargetDirectory, instance.Version, wrapper.WoUpdateOnly)
}

func (instance *Wrapper) executeEnsure(_ *kingpin.ParseContext) error {
	return wrapper.Write(instance.TargetDirectory, instance.Version, wrapper.WoCreateOrUpdate)
}
