package kubernetes

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func ResolveKubeConfig(contextName string) (*rest.Config, error) {
	clientConfig, err := rest.InClusterConfig()
	if err != nil {
		clientConfig, err = localConfig(contextName)
		if err != nil {
			return nil, err
		}
	}
	return clientConfig, nil
}

func localConfig(contextName string) (*rest.Config, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	apiConfig, err := rules.Load()
	if err != nil {
		return nil, err
	}
	if contextName != "" {
		apiConfig.CurrentContext = contextName
	}
	return clientcmd.NewDefaultClientConfig(*apiConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
}
