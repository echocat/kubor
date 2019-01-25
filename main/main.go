package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/levertonai/kubor/command"
	"github.com/levertonai/kubor/common"
	"github.com/levertonai/kubor/kubernetes"
	"github.com/levertonai/kubor/log"
	"github.com/levertonai/kubor/model"
	"github.com/levertonai/kubor/runtime"
	"os"
)

func version(_ *kingpin.ParseContext) error {
	_, err := fmt.Fprintln(os.Stderr, runtime.Runtime)
	return err
}

func main() {
	pf := model.NewProjectFactory()
	if err := command.Init(pf); err != nil {
		panic(err)
	}

	app := kingpin.New("kubor", "Safely bringing repositories using templating and charting inside CI/CD pipelines to Kubernetes.").
		ErrorWriter(os.Stderr).
		UsageWriter(os.Stderr).
		PreAction(func(_ *kingpin.ParseContext) error {
			return log.DefaultLogger.Init()
		})

	pf.ConfigureFlags(app)
	kubernetes.ConfigureKubeConfigFlags(app)
	log.DefaultLogger.ConfigureFlags(app)

	if err := common.ConfigureCliCommands("", app); err != nil {
		panic(err)
	}
	app.Command("version", "Print the actual version and other useful information.").
		Action(version)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
