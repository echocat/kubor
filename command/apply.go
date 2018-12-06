package command

import (
	"github.com/urfave/cli"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"kubor/common"
	"kubor/kubernetes"
	"kubor/model"
	"time"
)

func init() {
	cmd := &Apply{}
	cmd.Parent = cmd
	RegisterInitializable(cmd)
	common.RegisterCliFactory(cmd)
}

type Apply struct {
	Command

	Wait      time.Duration
	Predicate common.EvaluatingPredicate
}

func (instance *Apply) CreateCliCommands() ([]cli.Command, error) {
	return []cli.Command{{
		Name:   "apply",
		Usage:  "Apply the instances of this project using the provided values.",
		Action: instance.ExecuteFromCli,
		Flags: []cli.Flag{
			cli.DurationFlag{
				Name:        "wait, w",
				Usage:       "If set to value larger than 0 it will wait for this amount of time for successful running environment which was deployed. If it fails it will try to rollback.",
				EnvVar:      "KUBOR_WAIT",
				Value:       time.Minute * 5,
				Destination: &instance.Wait,
			},
			cli.GenericFlag{
				Name:   "predicate, p",
				Usage:  "Filters every object that should be listed. Empty allows everything. Pattern: \"[!]<template>=<must match regex>\", Example: \"{{.spec.name}}=Foo.*\"",
				EnvVar: "KUBOR_PREDICATE",
				Value:  &instance.Predicate,
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

	err = task.applySet.Execute()
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
