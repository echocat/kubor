package main

import (
	"encoding/base64"
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
	buildResources()
	buildBinaries(branch, commit)

	if withDocker {
		buildDockers(branch, false)
		tagDockers(branch)
	}
}

func buildResources() {
	buildWRapperResources()
}

func buildWRapperResources() {
	b := []byte(fmt.Sprintf(`
package wrapper

func init() {
	unixScript = "%s"
	windowsScript = "%s"
}
`, loadFileAsBase64("wrapper/kuborw"), loadFileAsBase64("wrapper/kuborw.cmd")))
	must(os.WriteFile("wrapper/resources_tmp.go", b, 0644))
}

func loadFile(source string) []byte {
	b, err := os.ReadFile(source)
	must(err)
	return b
}

func loadFileAsBase64(source string) string {
	b := loadFile(source)
	return base64.RawURLEncoding.EncodeToString(b)
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
		cmd.Env = append(os.Environ(), "GOOS="+t.os, "GOARCH="+t.arch, "CGO_ENABLED=0")
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

func buildDockers(branch string, forTesting bool) {
	prepareDockerResources()
	for _, v := range dockerVariants {
		buildDocker(branch, v, false, forTesting)
	}
}

func buildDocker(branch string, v dockerVariant, buildResources bool, forTesting bool) {
	if buildResources {
		prepareDockerResources()
	}
	version := branch
	if forTesting {
		version = "TEST" + version + "TEST"
	}
	execute(dockerCommand, "build", "-t", v.imageName(version), "-f", v.dockerFile, "--build-arg", "image="+imagePrefix, "--build-arg", "version="+version, ".")
}

func tagDockers(branch string) {
	for _, v := range dockerVariants {
		tagDocker(branch, v, false)
	}
}

func tagDocker(branch string, v dockerVariant, forTesting bool) {
	version := branch
	if forTesting {
		version = "TEST" + version + "TEST"
	}
	executeForVersionParts(version, func(tagSuffix string) {
		tagDockerWith(version, v, v.imageName(tagSuffix))
	})
	if latestVersionPattern != nil && latestVersionPattern.MatchString(version) {
		tagDockerWith(version, v, v.baseImageName())
	}
	if v.main {
		tagDockerWith(version, v, imagePrefix+":"+version)
		executeForVersionParts(version, func(tagSuffix string) {
			tagDockerWith(version, v, imagePrefix+":"+tagSuffix)
		})
		if latestVersionPattern != nil && latestVersionPattern.MatchString(version) {
			tagDockerWith(version, v, imagePrefix+":latest")
		}
	}
}

func tagDockerWith(branch string, v dockerVariant, tag string) {
	execute(dockerCommand, "tag", v.imageName(branch), tag)
}

func prepareDockerResources() {
	must(os.RemoveAll("var/docker/resources"))
	must(os.MkdirAll("var/docker/resources", 0755))
	download("https://dl.k8s.io/release/v"+kubectlVersion+"/bin/linux/amd64/kubectl", "var/docker/resources/usr/bin/kubectl", 0755)
	downloadFromTarGz("https://download.docker.com/linux/static/stable/x86_64/docker-"+dockerVersion+".tgz", "docker/docker", "var/docker/resources/usr/bin/docker", 0755)
	download("https://github.com/docker/machine/releases/download/v"+dockerMachineVersion+"/docker-machine-Linux-x86_64", "var/docker/resources/usr/bin/docker-machine", 0755)
}
