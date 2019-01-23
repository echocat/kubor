package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	_ = app.Command("build", "executes builds for the project").
		Action(func(*kingpin.ParseContext) error {
			build(branch, commit)
			return nil
		})
)

func build(branch, commit string) {
	buildBinaries(branch, commit)

	if withDocker {
		buildDockers(branch)
		tagDockers(branch)
	}
}

func buildBinaries(branch, commit string) {
	for _, t := range targets {
		buildBinary(branch, commit, t, false)
	}
}

func buildBinary(branch, commit string, t target, forTesting bool) {
	ldFlags := buildLdFlagsFor(branch, commit, forTesting)
	outputName := t.outputName()
	must(os.MkdirAll(filepath.Dir(outputName), 0755))
	executeTo(func(cmd *exec.Cmd) {
		cmd.Env = append(os.Environ(), "GOOS="+t.os, "GOARCH="+t.arch)
	}, os.Stderr, os.Stdout, "go", "build", "-ldflags", ldFlags, "-o", outputName, "./main")
}

func buildLdFlagsFor(branch, commit string, forTesting bool) string {
	testPrefix := ""
	testSuffix := ""
	if forTesting {
		testPrefix = "TEST"
		testSuffix = "TEST"
	}
	return fmt.Sprintf("-X main.extVersion=%s%s%s", testPrefix, branch, testSuffix) +
		fmt.Sprintf(" -X main.extRevision=%s%s%s", testPrefix, commit, testSuffix) +
		fmt.Sprintf(" -X main.extCompiled=%s", startTime.Format("2006-01-02T15:04:05Z"))
}

func buildDockers(branch string) {
	prepareDockerResources()
	for _, v := range dockerVariants {
		buildDocker(branch, v, false)
	}
}

func buildDocker(branch string, v dockerVariant, buildResources bool) {
	if buildResources {
		prepareDockerResources()
	}
	execute("docker", "build", "-t", v.imageName(branch), "-f", v.dockerFile, ".")
}

func tagDockers(branch string) {
	for _, v := range dockerVariants {
		tagDocker(branch, v)
	}
}

func tagDocker(branch string, v dockerVariant) {
	executeForVersionParts(branch, func(tagSuffix string) {
		tagDockerWith(branch, v, v.imageName(tagSuffix))
	})
	if latestVersionPattern != nil && latestVersionPattern.MatchString(branch) {
		tagDockerWith(branch, v, v.baseImageName())
	}
	if v.main {
		tagDockerWith(branch, v, imagePrefix+":"+branch)
		executeForVersionParts(branch, func(tagSuffix string) {
			tagDockerWith(branch, v, imagePrefix+":"+tagSuffix)
		})
		if latestVersionPattern != nil && latestVersionPattern.MatchString(branch) {
			tagDockerWith(branch, v, imagePrefix+":latest")
		}
	}
}

func tagDockerWith(branch string, v dockerVariant, tag string) {
	execute("docker", "tag", v.imageName(branch), tag)
}

func prepareDockerResources() {
	must(os.RemoveAll("var/docker/resources"))
	must(os.MkdirAll("var/docker/resources", 0755))
	download("https://storage.googleapis.com/kubernetes-release/release/v"+kubectlVersion+"/bin/linux/amd64/kubectl", "var/docker/resources/usr/bin/kubectl", 0755)
	downloadFromTarGz("https://download.docker.com/linux/static/stable/x86_64/docker-"+dockerVersion+".tgz", "docker/docker", "var/docker/resources/usr/bin/docker", 0755)
	download("https://github.com/docker/machine/releases/download/v"+dockerMachineVersion+"/docker-machine-Linux-x86_64", "var/docker/resources/usr/bin/docker-machine", 0755)
	download("https://github.com/theupdateframework/notary/releases/download/v"+dockerNotaryVersion+"/notary-Linux-amd64", "var/docker/resources/usr/bin/notary", 0755)
}