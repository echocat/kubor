package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/levertonai/kubor/command"
	"github.com/levertonai/kubor/common"
	"github.com/levertonai/kubor/kubernetes"
	"github.com/levertonai/kubor/log"
	"github.com/levertonai/kubor/model"
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
	extCompiled = ""
)

func version(_ *kingpin.ParseContext) error {
	_, err := fmt.Fprintf(os.Stderr, `kubor
 Version:      %s
 Built:        %s
 Git revision: %s
 Go version:   %s
 OS/Arch:      %s/%s
`,
		extVersion, extCompiled, extRevision, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	return err
}

func main() {
	if extCompiled == "" {
		extCompiled = time.Now().Format(timeFormat)
	}

	pf := model.NewProjectFactory()
	if err := command.Init(pf); err != nil {
		panic(err)
	}

	app := kingpin.New("kubor", "Safely bringing repositories using templating and charting inside CI/CD pipelines to Kubernetes.").
		Version(extVersion).
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
