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
	"runtime"
	"time"
)

const (
	timeFormat = "2006-01-02T15:04:05Z"
)

var (
	extVersion  = "development"
	extRevision = "development"
	extCompiled = time.Now().Format(timeFormat)

	versionCommand = cli.Command{
		Name:  "version",
		Usage: "Print the actual version and other useful information.",
		Action: func(*cli.Context) error {
			_, err := fmt.Fprintf(os.Stderr, `kubor
 Version:      %s
 Built:        %s
 Git revision: %s
 Go version:   %s
 OS/Arch:      %s/%s
`,
				extVersion, extCompiled, extRevision, runtime.Version(), runtime.GOOS, runtime.GOARCH)
			return err
		},
	}
)

func main() {
	var helpRequested bool
	cli.HelpFlag = cli.BoolFlag{
		Name:        "help, h",
		Usage:       "Show help",
		Destination: &helpRequested,
	}

	compiled, err := time.Parse(timeFormat, extCompiled)
	if err != nil {
		panic(fmt.Sprintf("illegal main.Compiled value '%s': %v", extCompiled, err))
	}

	pf := model.NewProjectFactory()
	err = command.Init(pf)
	if err != nil {
		panic(err)
	}
	app := cli.NewApp()
	app.Name = "kubor"
	app.Usage = "https://github.com/levertonai/kubor"
	app.Description = "Safely bringing repositories using templating and charting inside CI/CD pipelines to Kubernetes."
	app.Version = extVersion
	app.Compiled = compiled

	app.HideHelp = true
	app.HideVersion = true
	app.Flags = append(app.Flags, pf.Flags()...)
	app.Flags = append(app.Flags, kubernetes.KubeConfigFlags...)
	app.Flags = append(app.Flags, log.DefaultLogger.Flags()...)
	app.Flags = append(app.Flags, cli.HelpFlag)

	app.Commands, err = common.CreateCliCommands()
	if err != nil {
		panic(err)
	}
	app.Commands = append(app.Commands, versionCommand)

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
