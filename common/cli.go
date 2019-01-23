package common

import (
	"github.com/alecthomas/kingpin"
)

type HasFlags interface {
	Flag(name, help string) *kingpin.FlagClause
}

type HasCommands interface {
	Command(name, help string) *kingpin.CmdClause
}

type CliFactory interface {
	ConfigureCliCommands(HasCommands) error
}

var (
	cliFactories []CliFactory
)

func RegisterCliFactory(cliFactory CliFactory) CliFactory {
	cliFactories = append(cliFactories, cliFactory)
	return cliFactory
}

func ConfigureCliCommands(hc HasCommands) (err error) {
	for _, cliFactory := range cliFactories {
		if err := cliFactory.ConfigureCliCommands(hc); err != nil {
			return err
		}
	}
	return nil
}
