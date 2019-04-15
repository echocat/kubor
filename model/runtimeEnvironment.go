package model

import "io"

type RuntimeEnvironment struct {
	Module  Module            `json:"module"            yaml:"module"`
	Source  string            `json:"source,omitempty"  yaml:"source,omitempty"`
	Root    string            `json:"root,omitempty"    yaml:"root,omitempty"`
	Env     map[string]string `json:"env,omitempty"     yaml:"env,omitempty"`
	Context string            `json:"context,omitempty" yaml:"context,omitempty"`
}

func (instance *RuntimeEnvironment) Parse(source string, reader io.Reader) error {
	parser := RuntimeEnvironmentParser{
		Parser: DefaultStatementParser(),
	}
	return parser.Parse(source, instance, reader)
}
