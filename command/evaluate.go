package command

import (
	"fmt"
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

func (instance *Evaluate) ConfigureCliCommands(hc common.HasCommands) error {
	cmd := hc.Command("evaluate", "Evaluate the instances of this project using the provided values.").
		Alias("eval").
		Action(instance.ExecuteFromCli)

	cmd.Flag("sourceHint", "Prints to the output a comment which indicates where the rendered content organically comes from.").
		Envar("KUBOR_SOURCE_HINT").
		Default(fmt.Sprint(instance.SourceHint)).
		BoolVar(&instance.SourceHint)
	cmd.Flag("predicate", "Filters every object that should be listed. Empty allows everything."+
		" Example: \"{{.spec.name}}=Foo.*\"").
		PlaceHolder("[!]<template>=<must match regex>").
		Short('p').
		Envar("KUBOR_PREDICATE").
		SetValue(&instance.Predicate)

	return nil
}

func (instance *Evaluate) RunWithArguments(arguments Arguments) error {
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
