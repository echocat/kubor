package kubernetes

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/urfave/cli"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
	"kubor/log"
	"os"
	"path/filepath"
)

var (
	defaultKubeConfigPath = func() string {
		if home := homedir.HomeDir(); home != "" {
			return filepath.Join(home, ".kube", "config")
		}
		return ""
	}()
	kubeConfigPath  string
	kubeContext     string
	KubeConfigFlags = []cli.Flag{
		&cli.StringFlag{
			Name: "kubeconfig",
			Usage: "Path to the kubeconfig file. Optionally you can provide the content of the kubeconfig using\n" +
				"\tenvironment variable KUBE_CONFIG.",
			Value:       "",
			EnvVar:      "KUBOR_KUBECONFIG",
			Destination: &kubeConfigPath,
		},
		&cli.StringFlag{
			Name:        "context, c",
			Usage:       "Context of the kubeconfig which is used for the actual execution.",
			Value:       "",
			EnvVar:      "KUBOR_CONTEXT",
			Destination: &kubeContext,
		},
	}
)

func NewKubeConfig() (*restclient.Config, string, error) {
	clientConfig, contextName, err := NewKubeClientConfig()
	if err != nil {
		return nil, "", err
	}
	result, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, "", err
	}
	return result, contextName, nil
}

func NewKubeClientConfig() (clientcmd.ClientConfig, string, error) {
	selectedContext := kubeContext
	result := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&kubeConfigLoader{},
		&clientcmd.ConfigOverrides{
			CurrentContext: selectedContext,
		},
	)
	if selectedContext == "" {
		rc, err := result.RawConfig()
		if err != nil {
			return nil, "", err
		}
		selectedContext = rc.CurrentContext
	}
	log.
		WithField("context", selectedContext).
		Debug("Selected context: %v", selectedContext)
	return result, selectedContext, nil
}

type kubeConfigLoader struct {
	clientcmd.ClientConfigLoader
}

func (l *kubeConfigLoader) IsDefaultConfig(*restclient.Config) bool {
	return false
}

func (l *kubeConfigLoader) Load() (*clientcmdapi.Config, error) {
	config := clientcmdapi.NewConfig()
	atLeastOneConfigFound := false

	if plainFromEnv, ok := os.LookupEnv("KUBE_CONFIG"); ok {
		if fromEnv, err := clientcmd.Load([]byte(plainFromEnv)); err != nil {
			return nil, err
		} else if err := mergo.Merge(config, fromEnv); err != nil {
			return nil, err
		} else {
			atLeastOneConfigFound = true
		}
	}

	targetKubeConfigPath := kubeConfigPath
	if targetKubeConfigPath != "" {
		if _, err := os.Stat(targetKubeConfigPath); err != nil {
			return nil, err
		}
	} else {
		targetKubeConfigPath = defaultKubeConfigPath
	}

	if targetKubeConfigPath != "" {
		if fromFile, err := clientcmd.LoadFromFile(targetKubeConfigPath); os.IsNotExist(err) {
			// Ignore and continue
		} else if err != nil {
			return nil, err
		} else if err := mergo.Merge(config, fromFile); err != nil {
			return nil, err
		} else {
			atLeastOneConfigFound = true
		}
	}

	if !atLeastOneConfigFound {
		return nil, fmt.Errorf("there is neither argument --kubeconfig nor environment variable KUBE_CONFIG provided nor does %s exist", defaultKubeConfigPath)
	}

	return config, nil
}
