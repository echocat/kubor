package command

import (
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/tools/clientcmd"
	"kubor/common"
	"kubor/kubernetes"
	"os"
)

func init() {
	cmd := &Contexts{}
	common.RegisterCliFactory(cmd)
}

type Contexts struct{}

func (instance *Contexts) CreateCliCommands() ([]cli.Command, error) {
	return []cli.Command{{
		Name:  "contexts",
		Usage: "List available contexts",
		Action: func(*cli.Context) error {
			return instance.Run()
		},
	}}, nil
}

func (instance *Contexts) Run() error {
	config, currentContext, err := kubernetes.NewKubeClientConfig()
	if err != nil {
		return err
	}
	_, err = config.ClientConfig()
	if err != nil {
		return err
	}
	information, err := instance.toContextInformation(config, currentContext)
	if err != nil {
		return err
	}
	encoder := yaml.NewEncoder(os.Stdout)
	return encoder.Encode(information)
}

func (instance *Contexts) toContextInformation(config clientcmd.ClientConfig, currentContext string) ([]contextInformation, error) {
	rc, err := config.RawConfig()
	if err != nil {
		return nil, err
	}
	var result []contextInformation
	for name, context := range rc.Contexts {
		info := contextInformation{
			Name:     name,
			User:     context.AuthInfo,
			Selected: currentContext == name,
		}
		cluster := rc.Clusters[context.Cluster]
		if cluster != nil {
			info.Cluster = context.Cluster
			info.Server = cluster.Server
		}
		result = append(result, info)
	}
	return result, nil
}

type contextInformation struct {
	Name     string `yaml:"name,omitempty"`
	Selected bool   `yaml:"selected,omitempty"`
	Cluster  string `yaml:"cluster,omitempty"`
	Server   string `yaml:"server,omitempty"`
	User     string `yaml:"user,omitempty"`
}
