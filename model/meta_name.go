package model

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	metaNameDomainRegexp = regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	metaNameRegexp       = regexp.MustCompile(`^[a-z0-9A-Z]([a-z0-9A-Z._-]*[a-z0-9A-Z]|)$`)

	ErrIllegalMetaName = errors.New("illegal meta-name")
)

type MetaName string

func (instance *MetaName) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance MetaName) String() string {
	if v, err := instance.MarshalText(); err != nil {
		return fmt.Sprintf("illegal-meta-name-%s", string(instance))
	} else {
		return string(v)
	}
}

func (instance MetaName) MarshalText() (text []byte, err error) {
	parts := strings.SplitN(string(instance), "/", 2)
	domain, name := "", parts[0]
	if len(parts) > 1 {
		domain = parts[0]
		name = parts[1]

		if !metaNameDomainRegexp.MatchString(domain) || len(domain) > 253 {
			return nil, fmt.Errorf("%w: %s", ErrIllegalMetaName, string(instance))
		}
	}

	if name != "" && (!metaNameRegexp.MatchString(name) || len(name) > 63) {
		return nil, fmt.Errorf("%w: %s", ErrIllegalMetaName, string(instance))
	}

	return []byte(instance), nil
}

func (instance *MetaName) UnmarshalText(text []byte) error {
	result := MetaName(text)
	if _, err := result.MarshalText(); err != nil {
		return err
	}
	*instance = result
	return nil
}
