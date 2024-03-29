package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/echocat/slf4g"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	kubectlVersion       = "1.26.0"
	dockerVersion        = "23.0.1"
	dockerMachineVersion = "0.16.2"
)

var (
	versionRegexp = regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)$`)
	startTime     = time.Now()
)

func download(url string, to string, mode os.FileMode) {
	log.With("url", url).
		With("target", to).
		With("mode", mode).
		Info("Download...")
	body := startDownload(url)
	//noinspection GoUnhandledErrorResult
	defer body.Close()
	save(body, to, mode)
}

func downloadFromTarGz(url string, partName string, to string, mode os.FileMode) {
	log.With("url", url).
		With("part", partName).
		With("target", to).
		With("mode", mode).
		Info("Download...")

	body := startDownload(url)
	//noinspection GoUnhandledErrorResult
	defer body.Close()

	gr, err := gzip.NewReader(body)
	must(err)
	//noinspection GoUnhandledErrorResult
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		must(err)

		if header == nil {
			continue
		}

		if header.Name == partName {
			if header.Typeflag != tar.TypeReg {
				panic(fmt.Sprintf("Tar %s does cotnain the requested part '%s' but it is of unepxected type %d.", url, partName, header.Typeflag))
			}

			save(tr, to, mode)
			break
		}
	}
}

func startDownload(url string) io.ReadCloser {
	resp, err := http.Get(url)
	must(err)
	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("unexpected status while GET %s: %d - %s", url, resp.StatusCode, resp.Status))
	}
	return resp.Body
}

func save(from io.Reader, to string, mode os.FileMode) {
	dir := filepath.Dir(to)
	must(os.MkdirAll(dir, 0755))
	f, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	must(err)
	//noinspection GoUnhandledErrorResult
	defer f.Close()
	_, err = io.Copy(f, from)
	must(err)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func execute(args ...string) {
	executeTo(nil, os.Stderr, os.Stdout, args...)
}

func executeAndRecord(args ...string) string {
	buf := new(bytes.Buffer)
	executeTo(nil, buf, buf, args...)
	return buf.String()
}

type cmdCustomizer func(*exec.Cmd)

func executeTo(customizer cmdCustomizer, stderr, stdout io.Writer, args ...string) {
	if len(args) <= 0 {
		panic("no arguments provided")
	}
	log.With("command", args).
		Info("Execute...")

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	if customizer != nil {
		customizer(cmd)
	}

	if err := cmd.Run(); err != nil {
		msg := fmt.Sprintf("command failed [%s]: %v", strings.Join(args, " "), err)
		if b, ok := stdout.(fmt.Stringer); ok {
			msg += fmt.Sprintf("\nStdout: %s", b.String())
		}
		if b, ok := stderr.(fmt.Stringer); ok && stderr != stdout {
			msg += fmt.Sprintf("\nStderr: %s", b.String())
		}
		panic(msg)
	}
}

func quoteIfNeeded(what string) string {
	if strings.ContainsRune(what, '\t') ||
		strings.ContainsRune(what, '\n') ||
		strings.ContainsRune(what, ' ') ||
		strings.ContainsRune(what, '\xFF') ||
		strings.ContainsRune(what, '\u0100') ||
		strings.ContainsRune(what, '"') ||
		strings.ContainsRune(what, '\\') {
		return strconv.Quote(what)
	}
	return what
}

func quoteAllIfNeeded(in ...string) []string {
	out := make([]string, len(in))
	for i, a := range in {
		out[i] = quoteIfNeeded(a)
	}
	return out
}

func quoteAndJoin(in ...string) string {
	out := make([]string, len(in))
	for i, a := range in {
		out[i] = quoteIfNeeded(a)
	}
	return strings.Join(out, " ")
}

type versionPartAction func(versionPart string)

func executeForVersionParts(version string, action versionPartAction) {
	match := versionRegexp.FindStringSubmatch(version)
	if match != nil {
		action(fmt.Sprintf("v%s.%s", match[1], match[2]))
		action(fmt.Sprintf("v%s", match[1]))
	}
}
