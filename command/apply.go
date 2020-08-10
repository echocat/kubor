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

func init() {
	timeout := time.Minute * 5
	cmd := &Apply{
		Wait:       model.WaitUntil{Stage: model.WaitUntilStageApplied, Timeout: &timeout},
		KeepAlive:  1 * time.Minute,
		Predicate:  common.EvaluatingPredicate{},
		DryRun:     model.DryRunBefore,
		DryRunOn:   model.DryRunOnServerIfPossible,
		StageRange: model.StageRange{},
		Cleanup:    true,
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
	DryRun     model.DryRun
	DryRunOn   model.DryRunOn
	StageRange model.StageRange
	Cleanup    bool
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
		Default(instance.Wait.String()).
		SetValue(&instance.Wait)
	cmd.Flag("keepAlive", "If set to value larger than 0 it will do keep alive actions while wait for "+
		" completions.").
		Envar("KUBOR_KEEP_ALIVE").
		Default(instance.KeepAlive.String()).
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
	cmd.Flag("cleanup", "If enabled (default) it will remove all orphaned resources which matches the current"+
		" project's groupId and artifactId but where not part of the evaluated environment."+
		" This will be skipped in any way if either --predicate or --stage-range is defined.").
		Envar("KUBOR_CLEANUP").
		Default(fmt.Sprint(instance.Cleanup)).
		BoolVar(&instance.Cleanup)

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

func (instance *Apply) isCleanupAllowed() bool {
	return !instance.Predicate.IsRelevant() && !instance.StageRange.IsRelevant()
}

func (instance *Apply) RunWithArguments(arguments Arguments) error {
	ct, err := kubernetes.NewCleanupTask(arguments.Project, arguments.DynamicClient, kubernetes.CleanupModeOrphans)
	if err != nil {
		return err
	}
	task := &applyTask{
		source:        instance,
		dynamicClient: arguments.DynamicClient,
		arguments:     arguments,
		cleanupTask:   &ct,
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

	if instance.DryRun.IsDryRunAllowed() {
		if _, err := task.stagedApplySet.Execute("dryRun", instance.DryRunOn, nil, false); err != nil {
			return err
		}
	}

	if instance.DryRun.IsApplyAllowed() {
		if _, err := task.stagedApplySet.Execute("apply", model.DryRunNowhere, &instance.Wait, true); err != nil {
			return err
		}
	}

	if instance.Cleanup && instance.isCleanupAllowed() {
		if err := ct.Execute(); err != nil {
			return err
		}
	}

	return nil
}

type applyTask struct {
	source         *Apply
	dynamicClient  dynamic.Interface
	stagedApplySet kubernetes.StagedApplySet
	cleanupTask    *kubernetes.CleanupTask
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
	)
	if err != nil {
		return err
	}
	apply.KeepAliveInterval = instance.source.KeepAlive

	reference, err := kubernetes.GetObjectReference(object, instance.arguments.Project.Scheme)
	if err != nil {
		return err
	}

	stage, err := instance.arguments.Project.Annotations.GetStageFor(object)
	if err != nil {
		return err
	}

	if !instance.arguments.Project.Stages.Contains(stage) {
		return fmt.Errorf("%v (source: %s) has defined an unknown stage: %v; project defines: %v", reference, source, stage, instance.arguments.Project.Stages)
	}

	if err := instance.arguments.Project.Claim.Validate(reference); err != nil {
		return fmt.Errorf("%v (source: %s): %w", reference, source, err)
	}

	if instance.source.StageRange.Matches(instance.arguments.Project.Stages, stage) {
		instance.stagedApplySet.Add(stage, apply)
		instance.cleanupTask.Add(reference)
	}

	return nil
}
