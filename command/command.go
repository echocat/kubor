package command

import (
	"fmt"
	"github.com/urfave/cli"
	"k8s.io/client-go/dynamic"
	restclient "k8s.io/client-go/rest"
	"kubor/kubernetes"
	"kubor/model"
)

type CommandArguments struct {
	Project       model.Project
	Config        restclient.Config
	DynamicClient dynamic.Interface
}

type RunnableConsumingCommandArguments interface {
	RunWithArguments(args CommandArguments) error
}

type Command struct {
	ProjectFactory *model.ProjectFactory
	Parent         RunnableConsumingCommandArguments
}

func (instance *Command) Init(pf *model.ProjectFactory) error {
	instance.ProjectFactory = pf
	return nil
}

func (instance *Command) createProject(contextName string) (model.Project, error) {
	if instance.ProjectFactory == nil {
		return model.Project{}, fmt.Errorf("command not yet initialized")
	}
	return instance.ProjectFactory.Create(contextName)
}

func (instance *Command) ExecuteFromCli(*cli.Context) error {
	return instance.Run()
}

func (instance *Command) Run() error {
	config, contextName, err := kubernetes.NewKubeConfig()
	if err != nil {
		return err
	}
	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	project, err := instance.createProject(contextName)
	if err != nil {
		return err
	}
	if instance.Parent == nil {
		panic("no Parent defined")
	}
	return instance.Parent.RunWithArguments(CommandArguments{
		Project:       project,
		Config:        *config,
		DynamicClient: dc,
	})
}
