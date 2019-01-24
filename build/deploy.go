package main

import (
	"github.com/alecthomas/kingpin"
)

var (
	_ = app.Command("deploy", "executes deploys for the project").
		Action(func(*kingpin.ParseContext) error {
			deploy(branch)
			return nil
		})
	_ = app.Command("build-and-deploy", "executes builds and deploys for the project").
		Action(func(*kingpin.ParseContext) error {
			build(branch, commit)
			deploy(branch)
			return nil
		})
)

func deploy(branch string) {
	deployDockers(branch)
}

func deployDockers(branch string) {
	for _, v := range dockerVariants {
		deployDocker(branch, v)
	}
}

func deployDocker(branch string, v dockerVariant) {
	deployDockerTag(v.imageName(branch))
	executeForVersionParts(branch, func(tagSuffix string) {
		deployDockerTag(v.imageName(tagSuffix))
	})
	if latestVersionPattern != nil && latestVersionPattern.MatchString(branch) {
		deployDockerTag(v.baseImageName())
	}
	if v.main {
		deployDockerTag(imagePrefix + ":" + branch)
		executeForVersionParts(branch, func(tagSuffix string) {
			deployDockerTag(imagePrefix + ":" + tagSuffix)
		})
		if latestVersionPattern != nil && latestVersionPattern.MatchString(branch) {
			deployDockerTag(imagePrefix + ":latest")
		}
	}
}

func deployDockerTag(tag string) {
	execute("docker", "push", tag)
}
