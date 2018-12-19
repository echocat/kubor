package common

import (
	"github.com/urfave/cli"
)

type CliFactory interface {
	CreateCliCommands(context string) ([]cli.Command, error)
}

var (
	cliFactories []CliFactory
)

func RegisterCliFactory(cliFactory CliFactory) CliFactory {
	cliFactories = append(cliFactories, cliFactory)
	return cliFactory
}

func CreateCliCommands(context string) (result []cli.Command, err error) {
	for _, cliFactory := range cliFactories {
		if cc, err := cliFactory.CreateCliCommands(context); err != nil {
			return nil, err
		} else if cc != nil {
			result = append(result, cc...)
		}
	}
	return result, nil
}
