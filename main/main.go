package main

import (
	"fmt"
	"github.com/urfave/cli"
	"kubor/command"
	"kubor/common"
	"kubor/kubernetes"
	"kubor/log"
	"kubor/model"
	"os"
	"time"
)

var (
	extVersion  = "development"
	extCompiled = ""
)

func main() {
	pf := model.NewProjectFactory()
	err := command.Init(pf)
	if err != nil {
		panic(err)
	}
	cliCommands, err := common.CreateCliCommands()
	if err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Name = "kubor"
	app.Description = "Safely bringing repositories using templating and charting inside CI/CD pipelines to Kubernetes."
	app.Email = "info@leverton.ai"
	app.Version = extVersion
	//noinspection GoBoolExpressions
	if extCompiled != "" {
		compiled, err := time.Parse("2006-01-02T15:04:05Z", extCompiled)
		if err != nil {
			panic(fmt.Sprintf("illegal main.Compiled value '%s': %v", extCompiled, err))
		}
		app.Compiled = compiled
	}

	app.Flags = append(app.Flags, pf.Flags()...)
	app.Flags = append(app.Flags, kubernetes.KubeConfigFlags...)
	app.Flags = append(app.Flags, log.DefaultLogger.Flags()...)
	app.Commands = cliCommands
	app.Writer = os.Stderr
	app.ErrWriter = os.Stderr

	app.Before = func(*cli.Context) error {
		return log.DefaultLogger.Init()
	}

	if err := app.Run(os.Args); err != nil {
		//noinspection GoUnhandledErrorResult
		fmt.Fprint(app.ErrWriter, err.Error())
		os.Exit(1)
	}
}
