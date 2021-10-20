package models

import (
	"context"
	"fmt"
	radixclient "github.com/equinor/radix-operator/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	"k8s.io/client-go/rest"
	restclient "k8s.io/client-go/rest"
)

// Kube  Struct for accessing lower level kubernetes functions
type Kube struct {
	kubeClient  kubernetes.Interface
	radixClient radixclient.Interface
}

// GetKubeUtil Gets a Kube with kubernetes and Radix clients
func GetKubeUtil() (*Kube, error) {
	config, err := getKubernetesClientConfig()
	if err != nil {
		return nil, err
	}
	kubeClient, err := getKubernetesClientFromConfig(config)
	if err != nil {
		return nil, err
	}
	radixClient, err := radixclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Kube{kubeClient: kubeClient, radixClient: radixClient}, err
}

func getKubernetesClientConfig() (*restclient.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot get in-cluster config: %v", err)
	}
	return config, nil
}

func getKubernetesClientFromConfig(config *restclient.Config) (kubernetes.Interface, error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot create getClusterConfig k8s Radix client: %v", err)
	}
	return client, nil
}

//ExistsRadixRegistration Check if RadixRegistration exists by name
func (kubeUtil *Kube) ExistsRadixRegistration(name string) (bool, error) {
	rr, err := kubeUtil.radixClient.RadixV1().RadixRegistrations().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return rr != nil && rr.Name != "", err
}
