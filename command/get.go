package command

import (
	"github.com/levertonai/kubor/common"
	"github.com/levertonai/kubor/kubernetes"
	"github.com/levertonai/kubor/kubernetes/format"
	"github.com/levertonai/kubor/model"
	"k8s.io/client-go/dynamic"
	"os"
)

func init() {
	cmd := &Get{
		Output: format.VariantTable,
	}
	cmd.Parent = cmd
	RegisterInitializable(cmd)
	common.RegisterCliFactory(cmd)
}

type Get struct {
	Command

	Output    format.Variant
	Predicate common.EvaluatingPredicate
}

func (instance *Get) ConfigureCliCommands(context string, hc common.HasCommands) error {
	if context != "" {
		return nil
	}

	cmd := hc.Command("get", "Get the instances of this project using the provided values.").
		Action(instance.ExecuteFromCli)

	cmd.Flag("predicate", "Filters every object that should be listed. Empty allows everything. Pattern: \"[!]<template>=<must match regex>\", Example: \"{{.spec.name}}=Foo.*\"").
		Short('p').
		Envar("KUBOR_PREDICATE").
		SetValue(&instance.Predicate)
	cmd.Flag("output", "Defines how to format the output. Could be: table, yaml or json").
		Short('o').
		Default(instance.Output.String()).
		Envar("KUBOR_OUTPUT").
		SetValue(&instance.Output)

	return nil
}

func (instance *Get) RunWithArguments(arguments Arguments) error {
	if f, err := format.Get(instance.Output); err != nil {
		return err
	} else if ft, err := f.Format(os.Stdout); err != nil {
		return err
	} else {
		task := &getTask{
			source:        instance,
			dynamicClient: arguments.DynamicClient,
			claims:        &arguments.Project.Claims,
			formatTask:    ft,
		}
		defer ft.Close()

		for _, claim := range arguments.Project.Claims {
			for _, kind := range claim.Kinds {
				gvk := kind.ToGroupVersionKind()
				for _, namespace := range claim.Namespaces {
					if err := kubernetes.QueryNamespace(arguments.DynamicClient, gvk, namespace.String(), task.onObject); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
}

type getTask struct {
	source        *Get
	claims        *model.Claims
	dynamicClient dynamic.Interface
	formatTask    format.Task
}

func (instance *getTask) onObject(object kubernetes.Object) error {
	if matches, err := instance.claims.Matches(object); err != nil || !matches {
		return err
	}
	if matches, err := instance.source.Predicate.Matches(object); err != nil || !matches {
		return err
	}

	return instance.formatTask.Next(object)
}
