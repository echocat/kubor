package command

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/levertonai/kubor/kubernetes"
	"github.com/levertonai/kubor/model"
	"k8s.io/client-go/dynamic"
)

type Arguments struct {
	Project       model.Project
	Runtime       kubernetes.Runtime
	DynamicClient dynamic.Interface
}

type RunnableConsumingCommandArguments interface {
	RunWithArguments(args Arguments) error
}

type Command struct {
	ProjectFactory *model.ProjectFactory
	Parent         RunnableConsumingCommandArguments
}

func (instance *Command) Init(pf *model.ProjectFactory) error {
	instance.ProjectFactory = pf
	return nil
}

func (instance *Command) createProject(runtime kubernetes.Runtime) (model.Project, error) {
	if instance.ProjectFactory == nil {
		return model.Project{}, fmt.Errorf("command not yet initialized")
	}
	return instance.ProjectFactory.Create(runtime)
}

func (instance *Command) ExecuteFromCli(*kingpin.ParseContext) error {
	return instance.Run()
}

func (instance *Command) Run() error {
	runtime, err := kubernetes.NewRuntime()
	if err != nil {
		return err
	}
	dc, err := runtime.NewDynamicClient()
	if err != nil {
		return err
	}
	project, err := instance.createProject(runtime)
	if err != nil {
		return err
	}
	if instance.Parent == nil {
		panic("no Parent defined")
	}
	return instance.Parent.RunWithArguments(Arguments{
		Project:       project,
		Runtime:       runtime,
		DynamicClient: dc,
	})
}
