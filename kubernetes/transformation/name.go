package transformation

import (
	"errors"
	"fmt"
	"github.com/echocat/kubor/model"
	"regexp"
)

var (
	nameRegexp = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

	ErrIllegalName = errors.New("illegal transformation name")
)

type Name model.Name

func (instance *Name) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance Name) String() string {
	v, _ := instance.MarshalText()
	return string(v)
}

func (instance Name) MarshalText() (text []byte, err error) {
	if len(instance) > 0 && (!nameRegexp.MatchString(string(instance)) || len(instance) > 253) {
		return []byte(fmt.Sprintf("illegal-transformation-name-%s", string(instance))),
			fmt.Errorf("%w: %s", ErrIllegalName, string(instance))
	}
	return []byte(instance), nil
}

func (instance *Name) UnmarshalText(text []byte) error {
	v := Name(text)
	if _, err := instance.MarshalText(); err != nil {
		return err
	}
	*instance = v
	return nil
}
