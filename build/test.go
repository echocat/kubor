package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"os"
	"os/exec"
	"strings"
)

var (
	_ = app.Command("test", "executes tests for the project").
		Action(func(*kingpin.ParseContext) error {
			test(branch, commit)
			return nil
		})
)

func test(branch, commit string) {
	testGoCode(currentTarget)

	buildBinary(branch, commit, currentTarget, true)
	testBinary(branch, commit, currentTarget)

	if withDocker {
		buildBinary(branch, commit, linuxAmd64, true)
		for _, dv := range dockerVariants {
			buildDocker(branch, dv, true)
			testDocker(branch, commit, dv)
			tagDocker(branch, dv)
		}
	}
}

func testGoCode(t target) {
	executeTo(func(cmd *exec.Cmd) {
		cmd.Env = append(os.Environ(), "GOOS="+t.os, "GOARCH="+t.arch)
	}, os.Stderr, os.Stdout, "go", "test", "-v", "./...")
}

func testBinary(branch, commit string, t target) {
	testBinaryByExpectingResponse(t, `Version:      TEST`+branch+`TEST`, t.outputName(), "version")
	testBinaryByExpectingResponse(t, `Git revision: TEST`+commit+`TEST`, t.outputName(), "version")
}

func testBinaryByExpectingResponse(t target, expectedPartOfResponse string, args ...string) {
	cmd := append([]string{t.outputName()}, args...)
	response := executeAndRecord(args...)
	if !strings.Contains(response, expectedPartOfResponse) {
		panic(fmt.Sprintf("Command failed [%s]\nResponse should contain: %s\nBut response was: %s",
			quoteAllIfNeeded(cmd...), expectedPartOfResponse, response))
	}
}

func testDocker(branch, commit string, v dockerVariant) {
	testDockerByExpectingResponse(branch, v, "Version:      TEST"+branch+"TEST", "kubor", "version")
	testDockerByExpectingResponse(branch, v, "Git revision: TEST"+commit+"TEST", "kubor", "version")
	testDockerByExpectingResponse(branch, v, `GitVersion:"v`+kubectlVersion+`"`, "sh", "-c", "kubectl version || true")
	testDockerByExpectingResponse(branch, v, "Version:           "+dockerVersion+"\n", "sh", "-c", "docker version || true")
	testDockerByExpectingResponse(branch, v, "version "+dockerMachineVersion+",", "sh", "-c", "docker-machine version || true")
	testDockerByExpectingResponse(branch, v, "Version:    "+dockerNotaryVersion+"\n", "sh", "-c", "notary version || true")

}

func testDockerByExpectingResponse(branch string, v dockerVariant, expectedPartOfResponse string, command ...string) {
	call := append([]string{"docker", "run", "--rm", v.imageName(branch)}, command...)
	response := executeAndRecord(call...)
	if !strings.Contains(response, expectedPartOfResponse) {
		panic(fmt.Sprintf("Command failed [%s]\nResponse should contain: %s\nBut response was: %s",
			quoteAllIfNeeded(call...), expectedPartOfResponse, response))
	}
}