package main

import (
	"github.com/alecthomas/kingpin"
	"os"
	"regexp"
)

var (
	app = kingpin.New("build", "helps to build kubor").
		Interspersed(false)

	branch               = "snapshot"
	commit               = "unknown"
	withDocker           = true
	latestVersionPattern *regexp.Regexp
)

func init() {
	app.Flag("branch", "something like either main, v1.2.3 or snapshot-feature-foo").
		Required().
		Envar("TRAVIS_BRANCH").
		StringVar(&branch)
	app.Flag("commit", "something like 463e189796d5e96a7b605ab51985458faf8fd0d4").
		Required().
		Envar("TRAVIS_COMMIT").
		StringVar(&commit)
	app.Flag("withDocker", "enables docker tests and builds").
		Default("true").
		BoolVar(&withDocker)
	app.Flag("latestVersionPattern", "everything what matches here will be a latest tag").
		Envar("LATEST_VERSION_PATTERN").
		RegexpVar(&latestVersionPattern)
}

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
