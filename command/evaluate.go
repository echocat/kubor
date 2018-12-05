package command

import (
	"fmt"
	"github.com/urfave/cli"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
	"kubor/common"
	"kubor/model"
	"os"
)

func init() {
	cmd := &Evaluate{}
	cmd.Parent = cmd
	RegisterInitializable(cmd)
	common.RegisterCliFactory(cmd)
}

type Evaluate struct {
	Command

	SourceHint bool
	Predicate  common.EvaluatingPredicate
}

func (instance *Evaluate) CreateCliCommands() ([]cli.Command, error) {
	return []cli.Command{{
		Name:    "evaluate",
		Aliases: []string{"eval"},
		Usage:   "Evaluate the instances of this project using the provided values.",
		Action:  instance.ExecuteFromCli,
		Flags: []cli.Flag{
			cli.BoolTFlag{
				Name:        "sourceHint, sh",
				Usage:       "Prints to the output a comment which indicates where the rendered content organically comes from.",
				Destination: &instance.SourceHint,
			},
			cli.GenericFlag{
				Name:  "predicate, p",
				Usage: "Filters every object that should be listed. Empty allows everything. Pattern: \"[!]<template>=<must match regex>\", Example: \"{{.spec.name}}=Foo.*\"",
				Value: &instance.Predicate,
			},
		},
	}}, nil
}

func (instance *Evaluate) RunWithArguments(arguments CommandArguments) error {
	task := &evaluateTask{
		source: instance,
		first:  true,
	}
	oh, err := model.NewObjectHandler(task.onObject)
	if err != nil {
		return err
	}

	cp, err := arguments.Project.RenderedTemplatesProvider()
	if err != nil {
		return err
	}

	return oh.Handle(cp)
}

type evaluateTask struct {
	source *Evaluate
	first  bool
}

func (instance *evaluateTask) onObject(source string, object runtime.Object, unstructured *unstructured.Unstructured) error {
	if matches, err := instance.source.Predicate.Matches(unstructured.Object); err != nil {
		return err
	} else if !matches {
		return nil
	}

	if instance.first {
		instance.first = false
	} else {
		fmt.Print("---\n")
	}
	if instance.source.SourceHint {
		fmt.Printf(sourceHintTemplate, source)
	}
	encoder := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	return encoder.Encode(object, os.Stdout)
}
