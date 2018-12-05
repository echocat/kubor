package command

import (
	"fmt"
	"github.com/urfave/cli"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"kubor/common"
	"kubor/kubernetes"
	"kubor/model"
	"os"
)

func init() {
	cmd := &Get{}
	cmd.Parent = cmd
	RegisterInitializable(cmd)
	common.RegisterCliFactory(cmd)
}

type Get struct {
	Command

	Predicate  common.EvaluatingPredicate
	SourceHint bool
}

func (instance *Get) CreateCliCommands() ([]cli.Command, error) {
	return []cli.Command{{
		Name:   "get",
		Usage:  "Get the instances of this project using the provided values.",
		Action: instance.ExecuteFromCli,
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

func (instance *Get) RunWithArguments(arguments CommandArguments) error {
	task := &getTask{
		source:        instance,
		dynamicClient: arguments.DynamicClient,
		first:         true,
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

type getTask struct {
	source        *Get
	dynamicClient dynamic.Interface
	first         bool
}

func (instance *getTask) onObject(source string, object runtime.Object, unstructured *unstructured.Unstructured) error {
	if matches, err := instance.source.Predicate.Matches(unstructured.Object); err != nil {
		return err
	} else if !matches {
		return nil
	}

	resource, err := kubernetes.GetObjectResource(unstructured, instance.dynamicClient)
	if err != nil {
		return err
	}
	ul, err := resource.Get(nil)
	if err != nil {
		return err
	}

	if instance.first {
		instance.first = false
	} else {
		fmt.Print("---\n")
	}
	if instance.source.SourceHint {
		fmt.Printf(sourceHintTemplate, source)
	}

	if err := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme).
		Encode(ul, os.Stdout); err != nil {
		return err
	}

	return nil
}
