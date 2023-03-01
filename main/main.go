package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/echocat/kubor/command"
	"github.com/echocat/kubor/common"
	"github.com/echocat/kubor/kubernetes"
	"github.com/echocat/kubor/model"
	"github.com/echocat/slf4g/native"
	"github.com/echocat/slf4g/native/facade/value"
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
		ErrorWriter(os.Stderr).
		UsageWriter(os.Stderr)

	lv := value.NewProvider(native.DefaultProvider)

	app.Flag("logLevel", "Specifies the minimum required log level.").
		Envar("KUBOR_LOG_LEVEL").
		Default(lv.Level.String()).
		SetValue(&lv.Level)
	app.Flag("logFormat", "Specifies format output (text or json).").
		Envar("KUBOR_LOG_FORMAT").
		Default(lv.Consumer.Formatter.String()).
		SetValue(&lv.Consumer.Formatter)
	app.Flag("logColorMode", "Specifies if the output is in colors or not (auto, never or always).").
		Envar("KUBOR_LOG_COLOR_MODE").
		Default(lv.Consumer.Formatter.ColorMode.String()).
		SetValue(lv.Consumer.Formatter.ColorMode)

	pf.ConfigureFlags(app)
	kubernetes.ConfigureKubeConfigFlags(app)

	if err := common.ConfigureCliCommands("", app, extVersion); err != nil {
		panic(err)
	}
	app.Command("version", "Print the actual version and other useful information.").
		Action(version)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
