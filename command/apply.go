package command

import (
	"fmt"
	"github.com/urfave/cli"
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
		DryRun: DryRunType("before"),
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
}

func (instance *Apply) CreateCliCommands(context string) ([]cli.Command, error) {
	if context != "" {
		return nil, nil
	}
	return []cli.Command{{
		Name:   "apply",
		Usage:  "Apply the instances of this project using the provided values.",
		Action: instance.ExecuteFromCli,
		Flags: []cli.Flag{
			cli.DurationFlag{
				Name: "wait, w",
				Usage: "If set to value larger than 0 it will wait for this amount of time for successful\n" +
					"\trunning environment which was deployed. If it fails it will try to rollback.",
				EnvVar:      "KUBOR_WAIT",
				Value:       time.Minute * 5,
				Destination: &instance.Wait,
			},
			cli.GenericFlag{
				Name: "predicate, p",
				Usage: "Filters every object that should be listed. Empty allows everything.\n" +
					"\tPattern: \"[!]<template>=<must match regex>\", Example: \"{{.spec.name}}=Foo.*\"",
				EnvVar: "KUBOR_PREDICATE",
				Value:  &instance.Predicate,
			},
			cli.GenericFlag{
				Name: "dryRun",
				Usage: "If set to 'before' it will execute a dry run before the actual apply.\n" +
					"\tThis is perfect in cases where the first parts of the apply configuration works and\n" +
					"\tthe following stuff is broken. If set to 'never' apply will be executed without dry run.\n" +
					"\tOn 'only' it will only run the dry run but not the apply.",
				EnvVar: "KUBOR_DRY_RUN",
				Value:  &instance.DryRun,
			},
		},
	}}, nil
}

func (instance *Apply) RunWithArguments(arguments CommandArguments) error {
	task := &applyTask{
		source:        instance,
		dynamicClient: arguments.DynamicClient,
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
		err = task.applySet.Execute(true)
		if err != nil {
			return err
		}
		if instance.DryRun == DryRunType("only") {
			return nil
		}
	}

	err = task.applySet.Execute(false)
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
}

func (instance *applyTask) onObject(source string, object runtime.Object, unstructured *unstructured.Unstructured) error {
	if matches, err := instance.source.Predicate.Matches(unstructured.Object); err != nil {
		return err
	} else if !matches {
		return nil
	}

	apply, err := kubernetes.NewApplyObject(source, unstructured, instance.dynamicClient)
	if err != nil {
		return err
	}

	instance.applySet.Add(apply)

	return nil
}
