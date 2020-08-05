package command

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/echocat/kubor/common"
	"github.com/echocat/kubor/kubernetes"
	"github.com/echocat/kubor/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"time"
)

type DryRunType string

const (
	DryRunBefore = DryRunType("before")
	DryRunNever  = DryRunType("never")
	DryRunOnly   = DryRunType("only")
)

func (instance *DryRunType) Set(plain string) error {
	candidate := DryRunType(plain)
	switch candidate {
	case DryRunBefore, DryRunNever, DryRunOnly:
		*instance = candidate
		return nil
	default:
		return fmt.Errorf("unsupported dry run type: %s", plain)
	}
}

func (instance DryRunType) String() string {
	return string(instance)
}

func init() {
	timeout := time.Minute * 5
	cmd := &Apply{
		Wait:     model.WaitUntil{Stage: model.WaitUntilStageApplied, Timeout: &timeout},
		DryRun:   DryRunType("before"),
		DryRunOn: kubernetes.ServerIfPossibleDryRun,
	}
	cmd.Parent = cmd
	RegisterInitializable(cmd)
	common.RegisterCliFactory(cmd)
}

type Apply struct {
	Command

	Wait       model.WaitUntil
	KeepAlive  time.Duration
	Predicate  common.EvaluatingPredicate
	DryRun     DryRunType
	DryRunOn   kubernetes.DryRunOn
	StageRange model.StageRange
}

func (instance *Apply) ConfigureCliCommands(context string, hc common.HasCommands, _ string) error {
	if context != "" {
		return nil
	}

	cmd := hc.Command("apply", "Apply the instances of this project using the provided values.").
		Action(instance.ExecuteFromCli)

	cmd.Flag("wait", "If set to value larger than 0 it will wait for this amount of time for successful"+
		" running environment which was deployed. If it fails it will try to rollback.").
		Short('w').
		Envar("KUBOR_WAIT").
		Default("applied:" + (time.Minute * 5).String()).
		SetValue(&instance.Wait)
	cmd.Flag("keepAlive", "If set to value larger than 0 it will do keep alive actions while wait for "+
		" completions.").
		Envar("KUBOR_KEEP_ALIVE").
		Default((time.Minute * 1).String()).
		DurationVar(&instance.KeepAlive)
	cmd.Flag("predicate", "Filters every object that should be listed. Empty allows everything."+
		" Example: \"{{.spec.name}}=Foo.*\"").
		PlaceHolder("[!]<template>=<must match regex>").
		Short('p').
		Envar("KUBOR_PREDICATE").
		SetValue(&instance.Predicate)
	cmd.Flag("dryRun", "If set to 'before' it will execute a dry run before the actual apply."+
		" This is perfect in cases where the first parts of the apply configuration works and"+
		" the following stuff is broken. If set to 'never' apply will be executed without dry run."+
		" On 'only' it will only run the dry run but not the apply.").
		Envar("KUBOR_DRY_RUN").
		Default(instance.DryRun.String()).
		SetValue(&instance.DryRun)
	cmd.Flag("dryRunOn", "If set to 'server' it will execute the dry run on the target kubernetes server"+
		" if this is not supported the apply will fail."+
		" If set to 'client' it will only run inside kubor and never will call the server at all."+
		" If set to 'serverIfPossible' it will check if it is available to run on the server if not it will just run"+
		" inside kubor.").
		Envar("KUBOR_DRY_RUN_ON").
		Default(instance.DryRunOn.String()).
		SetValue(&instance.DryRunOn)
	cmd.Flag("stageRange", "If set it will specify from which to which stage kubor will execute"+
		" the deployment."+
		" Pattern: [<from-stage>]:[<to-stage>]. If one of the terms it means not limited.").
		Envar("KUBOR_STAGE_RANGE").
		Default(instance.StageRange.String()).
		SetValue(&instance.StageRange)

	cmd.Validate(func(clause *kingpin.CmdClause) error {
		switch instance.Wait.Stage {
		case model.WaitUntilStageApplied, model.WaitUntilStageNever:
			return nil
		default:
			return fmt.Errorf("--wait only support 'applied' or 'never', but got: %v", instance.Wait.Stage)
		}
	})

	return nil
}

func (instance *Apply) RunWithArguments(arguments Arguments) error {
	task := &applyTask{
		source:        instance,
		dynamicClient: arguments.DynamicClient,
		arguments:     arguments,
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

	if instance.DryRun == DryRunBefore || instance.DryRun == DryRunOnly {
		err = task.stagedApplySet.DryRun(instance.DryRunOn)
		if err != nil {
			return err
		}
		if instance.DryRun == DryRunOnly {
			return nil
		}
	}

	_, err = task.stagedApplySet.Execute(instance.Wait)
	return err
}

type applyTask struct {
	source         *Apply
	dynamicClient  dynamic.Interface
	stagedApplySet kubernetes.StagedApplySet
	arguments      Arguments
}

func (instance *applyTask) onObject(source string, _ runtime.Object, object *unstructured.Unstructured) error {
	if matches, err := instance.source.Predicate.Matches(object.Object); err != nil {
		return err
	} else if !matches {
		return nil
	}

	apply, err := kubernetes.NewApplyObject(
		instance.arguments.Project,
		source,
		object,
		instance.dynamicClient,
		instance.arguments.Runtime,
		instance.arguments.Project.Validation.Schema,
	)
	if err != nil {
		return err
	}
	apply.KeepAliveInterval = instance.source.KeepAlive

	stage, err := instance.arguments.Project.Annotations.GetStageFor(object)
	if err != nil {
		return err
	}

	if instance.source.StageRange.Matches(instance.arguments.Project.Stages, stage) {
		instance.stagedApplySet.Add(stage, apply)
	}

	return nil
}
