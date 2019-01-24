package command

import (
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"kubor/common"
	"kubor/kubernetes"
	"kubor/model"
	"time"
)

type DryRunType string

func (instance *DryRunType) Set(plain string) error {
	if plain != "before" && plain != "never" && plain != "only" {
		return fmt.Errorf("unsupported dry run type: %s", plain)
	}
	*instance = DryRunType(plain)
	return nil
}

func (instance DryRunType) String() string {
	return string(instance)
}

func init() {
	cmd := &Apply{
		DryRun:   DryRunType("before"),
		DryRunOn: kubernetes.ServerIfPossibleDryRun,
	}
	cmd.Parent = cmd
	RegisterInitializable(cmd)
	common.RegisterCliFactory(cmd)
}

type Apply struct {
	Command

	Wait      time.Duration
	Predicate common.EvaluatingPredicate
	DryRun    DryRunType
	DryRunOn  kubernetes.DryRunOn
}

func (instance *Apply) ConfigureCliCommands(hc common.HasCommands) error {
	cmd := hc.Command("apply", "Apply the instances of this project using the provided values.").
		Action(instance.ExecuteFromCli)

	cmd.Flag("wait", "If set to value larger than 0 it will wait for this amount of time for successful"+
		" running environment which was deployed. If it fails it will try to rollback.").
		Short('w').
		Envar("KUBOR_WAIT").
		Default((time.Minute * 5).String()).
		DurationVar(&instance.Wait)
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

	return nil
}

func (instance *Apply) RunWithArguments(arguments Arguments) error {
	task := &applyTask{
		source:        instance,
		dynamicClient: arguments.DynamicClient,
		runtime:       arguments.Runtime,
	}
	oh, err := model.NewObjectHandler(task.onObject)
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

	if instance.DryRun == DryRunType("before") || instance.DryRun == DryRunType("only") {
		err = task.applySet.Execute(instance.DryRunOn)
		if err != nil {
			return err
		}
		if instance.DryRun == DryRunType("only") {
			return nil
		}
	}

	err = task.applySet.Execute(kubernetes.NowhereDryRun)
	if err != nil {
		return err
	}

	if task.source.Wait <= 0 {
		return nil
	}

	return task.applySet.Wait(task.source.Wait)
}

type applyTask struct {
	source        *Apply
	dynamicClient dynamic.Interface
	applySet      kubernetes.ApplySet
	runtime       kubernetes.Runtime
}

func (instance *applyTask) onObject(source string, object runtime.Object, unstructured *unstructured.Unstructured) error {
	if matches, err := instance.source.Predicate.Matches(unstructured.Object); err != nil {
		return err
	} else if !matches {
		return nil
	}

	apply, err := kubernetes.NewApplyObject(source, unstructured, instance.dynamicClient, instance.runtime)
	if err != nil {
		return err
	}

	instance.applySet.Add(apply)

	return nil
}
