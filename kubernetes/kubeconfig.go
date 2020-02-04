package kubernetes

import (
	"errors"
	"fmt"
	"github.com/echocat/kubor/common"
	"github.com/echocat/kubor/log"
	"github.com/imdario/mergo"
	"k8s.io/client-go/discovery"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	defaultKubeConfigPath = func() string {
		if home := homedir.HomeDir(); home != "" {
			return filepath.Join(home, ".kube", "config")
		}
		return ""
	}()
	kubeConfigPath string
	kubeContext    string
)

func ConfigureKubeConfigFlags(hf common.HasFlags) {
	hf.Flag("kubeconfig", "Path to the kubeconfig file. Optionally you can provide the content of the kubeconfig using"+
		" environment variable KUBE_CONFIG."+
		" If this value is 'mock' it will use a kubeconfig which cannot do deployments but works with every named context"+
		"; --context is in this case mandatory.").
		Envar("KUBECONFIG").
		PlaceHolder("<kube config file>").
		StringVar(&kubeConfigPath)
	hf.Flag("context", "Context of the kubeconfig which is used for the actual execution.").
		Short('c').
		Envar("KUBOR_CONTEXT").
		PlaceHolder("<context>").
		StringVar(&kubeContext)
}

func NewRuntime() (Runtime, error) {
	if kubeConfigPath == "mock" {
		if kubeContext == "" {
			kubeContext = "mock"
		}
		return newRuntimeMock(kubeContext)
	}
	clientConfig, contextName, err := NewKubeClientConfig()
	if err != nil {
		return nil, err
	}
	return newRuntimeImpl(clientConfig, contextName)
}

func NewKubeClientConfig() (clientcmd.ClientConfig, string, error) {
	if kubeConfigPath == "mock" {
		return nil, "", errors.New("this operation is not supported if --kubeconfig=mock was specified")
	}
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

func newDiscoveryClientFor(config *restclient.Config) (*discovery.CachedDiscoveryClient, error) {
	discoveryCacheDir := computeDiscoverCacheDir(filepath.Join(homedir.HomeDir(), ".kube", "cache", "discovery"), config.Host)
	httpCacheDir := filepath.Join(homedir.HomeDir(), ".kube", "http-cache")
	return discovery.NewCachedDiscoveryClientForConfig(config, discoveryCacheDir, httpCacheDir, time.Duration(10*time.Minute))
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

var overlyCautiousIllegalFileCharacters = regexp.MustCompile(`[^(\w/.)]`)

func computeDiscoverCacheDir(parentDir, host string) string {
	// strip the optional scheme from host if its there:
	schemelessHost := strings.Replace(strings.Replace(host, "https://", "", 1), "http://", "", 1)
	// now do a simple collapse of non-AZ09 characters.  Collisions are possible but unlikely.  Even if we do collide the problem is short lived
	safeHost := overlyCautiousIllegalFileCharacters.ReplaceAllString(schemelessHost, "_")
	return filepath.Join(parentDir, safeHost)
}
