package wrapper

import (
	"fmt"
	"strings"
)

type WriteOpt string

const (
	WoCreateOrUpdate = "createOrUpdate"
	WoCreateOnly     = "createOnly"
	WoUpdateOnly     = "updateOnly"
)

var WriteOpts = []WriteOpt{WoCreateOrUpdate, WoCreateOnly, WoUpdateOnly}

func (instance WriteOpt) String() string {
	return string(instance)
}

func (instance *WriteOpt) Set(plain string) error {
	for _, candidate := range WriteOpts {
		if strings.ToLower(string(candidate)) == strings.ToLower(plain) {
			*instance = candidate
		}
	}
	return fmt.Errorf("unknown write option: %s", plain)
}

//func Write(targetDir string, version string, opt WriteOpt) error {
//
//}
//
//func writeFile(target string, base64Content string, version string, opt WriteOpt) error {
//	if b, err := base64.RawURLEncoding.DecodeString(base64Content); err != nil {
//		return err
//	} else if f, err := {
//		b = []byte(strings.Replace(string(b), "####UNDEFINED####", version, -1))
//
//	}
//	if
//
//}
