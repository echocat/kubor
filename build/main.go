package main

import (
	"github.com/alecthomas/kingpin"
	_ "github.com/echocat/slf4g/native"
	"os"
	"regexp"
)

var (
	app = kingpin.New("build", "helps to build kubor").
		Interspersed(false)

	branch               = "snapshot"
	commit               = "unknown"
	withDocker           = true
	dockerCommand        = "docker"
	latestVersionPattern *regexp.Regexp
)

func init() {
	app.Flag("branch", "something like either main, v1.2.3 or snapshot-feature-foo").
		Required().
		Envar("GITHUB_REF_NAME").
		StringVar(&branch)
	app.Flag("commit", "something like 463e189796d5e96a7b605ab51985458faf8fd0d4").
		Required().
		Envar("GITHUB_SHA").
		StringVar(&commit)
	app.Flag("docker.enabled", "enables docker tests and builds").
		Default("true").
		BoolVar(&withDocker)
	app.Flag("docker.command", "command to use on docker builds").
		Default("docker").
		StringVar(&dockerCommand)
	app.Flag("latestVersionPattern", "everything what matches here will be a latest tag").
		Envar("LATEST_VERSION_PATTERN").
		RegexpVar(&latestVersionPattern)
}

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
