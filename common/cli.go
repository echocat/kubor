package common

import (
	"github.com/urfave/cli"
)

type CliFactory interface {
	CreateCliCommands() ([]cli.Command, error)
}

var (
	cliFactories []CliFactory
)

func RegisterCliFactory(cliFactory CliFactory) CliFactory {
	cliFactories = append(cliFactories, cliFactory)
	return cliFactory
}

func CreateCliCommands() (result []cli.Command, err error) {
	for _, cliFactory := range cliFactories {
		if cc, err := cliFactory.CreateCliCommands(); err != nil {
			return nil, err
		} else {
			result = append(result, cc...)
		}
	}
	return result, nil
}
