package main

import (
	"fmt"
	"path/filepath"
	"runtime"
)

const imagePrefix = "echocat/kubor"

var (
	dockerVariants = []dockerVariant{
		{
			base:       "alpine",
			dockerFile: "Dockerfile",
			main:       true,
		},
		{
			base:       "ubuntu",
			dockerFile: "Dockerfile.ubuntu",
		},
	}

	currentTarget = target{os: runtime.GOOS, arch: runtime.GOARCH}
	linuxAmd64    = target{os: "linux", arch: "amd64"}
	targets       = []target{
		{os: "darwin", arch: "amd64"},
		{os: "darwin", arch: "386"},
		linuxAmd64,
		{os: "linux", arch: "386"},
		{os: "windows", arch: "amd64"},
		{os: "windows", arch: "386"},
	}
)

type dockerVariant struct {
	base       string
	dockerFile string
	main       bool
}

func (instance dockerVariant) baseImageName() string {
	return imagePrefix + ":" + instance.base
}

func (instance dockerVariant) imageName(branch string) string {
	result := imagePrefix + ":" + instance.base
	if branch != "" {
		result += "-" + branch
	}
	return result
}

type target struct {
	os   string
	arch string
}

func (instance target) outputName() string {
	return filepath.Join("dist", fmt.Sprintf("kubor-%s-%s%s", instance.os, instance.arch, instance.ext()))
}

func (instance target) ext() string {
	if instance.os == "windows" {
		return ".exe"
	}
	return ""
}
