package command

import (
	"fmt"
	"github.com/echocat/kubor/common"
	"github.com/echocat/kubor/kubernetes"
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func init() {
	cmd := &Cleanup{}
	cmd.Parent = cmd
	RegisterInitializable(cmd)
	common.RegisterCliFactory(cmd)
}

type Cleanup struct {
	Command
}

func (instance *Cleanup) ConfigureCliCommands(context string, hc common.HasCommands, _ string) error {
	if context != "" {
		return nil
	}

	hc.Command("cleanup", "Will remove all orphaned resources which matches the current"+
		" project's groupId and artifactId but where not part of the evaluated environment.").
		Action(instance.ExecuteFromCli)
	return nil
}

func (instance *Cleanup) RunWithArguments(arguments Arguments) error {
	ct, err := kubernetes.NewCleanupTask(arguments.Project, arguments.DynamicClient)
	if err != nil {
		return err
	}
	task := &cleanupTask{
		source:      instance,
		project:     arguments.Project,
		cleanupTask: &ct,
	}
	oh, err := model.NewObjectHandler(task.onObject, arguments.Project)
	if err != nil {
		return err
	}

	cp, err := arguments.Project.RenderedTemplatesProvider()
	if err != nil {
		return err
	}

	err = oh.Handle(cp)
	if err != nil {
		return err
	}

	return ct.Execute()
}

type cleanupTask struct {
	source      *Cleanup
	project     *model.Project
	cleanupTask *kubernetes.CleanupTask
}

func (instance *cleanupTask) onObject(source string, _ runtime.Object, object *unstructured.Unstructured) error {
	reference, err := kubernetes.GetObjectReference(object, instance.project.Scheme)
	if err != nil {
		return err
	}

	if err := instance.project.Claim.Validate(reference); err != nil {
		return fmt.Errorf("%v (source: %s): %w", reference, source, err)
	}

	instance.cleanupTask.Add(reference)

	return nil
}
