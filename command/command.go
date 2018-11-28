package command

import (
	"fmt"
	"github.com/urfave/cli"
	"k8s.io/client-go/dynamic"
	restclient "k8s.io/client-go/rest"
	"kubor/kubernetes"
	"kubor/model"
)

type RunnableConsumingProject interface {
	RunForProject(project model.Project) error
}

type Command struct {
	ProjectFactory *model.ProjectFactory
	Parent         RunnableConsumingProject
}

func (instance *Command) Init(pf *model.ProjectFactory) error {
	instance.ProjectFactory = pf
	return nil
}

func (instance *Command) createProject() (model.Project, error) {
	if instance.ProjectFactory == nil {
		return model.Project{}, fmt.Errorf("command not yet initialized")
	}
	return instance.ProjectFactory.Create()
}

func (instance *Command) ExecuteFromCli(*cli.Context) error {
	return instance.Run()
}

func (instance *Command) Run() error {
	project, err := instance.createProject()
	if err != nil {
		return err
	}
	if instance.Parent == nil {
		panic("no Parent defined")
	}
	return instance.Parent.RunForProject(project)
}

type RunnableConsumingKubernetesClientCommandArguments interface {
	RunForArguments(arguments KubernetesClientCommandArguments) error
}

type KubernetesClientCommand struct {
	Command
	Parent RunnableConsumingKubernetesClientCommandArguments
}

type KubernetesClientCommandArguments struct {
	Project       model.Project
	Config        restclient.Config
	DynamicClient dynamic.Interface
}

func (instance *KubernetesClientCommand) Init(pf *model.ProjectFactory) error {
	if err := instance.Command.Init(pf); err != nil {
		return err
	}
	instance.Command.Parent = instance
	return nil
}

func (instance *KubernetesClientCommand) RunForProject(project model.Project) error {
	config, err := kubernetes.NewKubeConfig()
	if err != nil {
		return err
	}
	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	if instance.Parent == nil {
		panic("no Parent defined")
	}
	return instance.Parent.RunForArguments(KubernetesClientCommandArguments{
		Project:       project,
		Config:        *config,
		DynamicClient: dc,
	})
}
