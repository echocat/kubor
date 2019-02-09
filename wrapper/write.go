package wrapper

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
)

type WriteOpt string

const (
	WoCreateOrUpdate = "createOrUpdate"
	WoCreateOnly     = "createOnly"
	WoUpdateOnly     = "updateOnly"
)

var (
	unixScript    = ``
	windowsScript = ``
)

func Write(targetDir string, version string, opt WriteOpt) error {
	if err := writeFile(filepath.Join(targetDir, "kuborw"), unixScript, version, opt, 0755); err != nil {
		return err
	} else if err := writeFile(filepath.Join(targetDir, "kuborw.cmd"), windowsScript, version, opt, 0644); err != nil {
		return err
	} else {
		return nil
	}
}

func writeFile(target string, rawBase64EncodedContent string, version string, opt WriteOpt, perm os.FileMode) error {
	if content, err := prepareContent(rawBase64EncodedContent, version); err != nil {
		return err
	} else if err := createDirectorsForFileIfRequired(target, opt); err != nil {
		return err
	} else if f, err := openFile(target, opt, perm); err != nil {
		return err
	} else {
		defer f.Close()
		_, err := f.Write(content)
		return err
	}
}

func prepareContent(rawBase64EncodedContent string, version string) ([]byte, error) {
	if b, err := base64.RawURLEncoding.DecodeString(rawBase64EncodedContent); err != nil {
		return nil, err
	} else {
		return []byte(strings.Replace(string(b), "####VERSION####", version, -1)), nil
	}
}

func createDirectorsForFileIfRequired(file string, opt WriteOpt) error {
	if opt == WoUpdateOnly {
		return nil
	}
	return os.MkdirAll(filepath.Dir(file), 0755)
}

func openFile(file string, opt WriteOpt, perm os.FileMode) (*os.File, error) {
	of := os.O_WRONLY
	switch opt {
	case WoUpdateOnly:
		of |= os.O_TRUNC
	case WoCreateOnly:
		of |= os.O_CREATE | os.O_EXCL
	default:
		of |= os.O_CREATE | os.O_TRUNC
	}
	return os.OpenFile(file, of, perm)
}
