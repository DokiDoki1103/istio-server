package kubernetes

import (
	"context"
	kube "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"

	istio "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/rest"
	gatewayapiclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

var (
	k8sClient *K8SClient
)

func GetK8sClient() *K8SClient {
	return k8sClient
}

type K8SClient struct {
	configPath     string // kubernetes config file path
	k8s            kube.Interface
	istioClientset istio.Interface
	restConfig     *rest.Config
	ctx            context.Context
	isOpenShift    *bool
	isGatewayAPI   *bool
	gatewayapi     gatewayapiclient.Interface
	isIstioAPI     *bool
}

func NewClientFromConfig() (*K8SClient, error) {
	client := K8SClient{
		configPath: os.Getenv("K8S_CONFIG"),
	}

	if client.configPath == "" {
		config, err := restclient.InClusterConfig()
		if err != nil {
			return nil, err
		}
		client.restConfig = config
	} else {
		config, err := clientcmd.BuildConfigFromFlags("", client.configPath)
		if err != nil {
			return nil, err
		}
		client.restConfig = config
	}

	k8s, err := kube.NewForConfig(client.restConfig)
	if err != nil {
		return nil, err
	}
	client.k8s = k8s

	client.istioClientset, err = istio.NewForConfig(client.restConfig)
	if err != nil {
		return nil, err
	}

	client.gatewayapi, err = gatewayapiclient.NewForConfig(client.restConfig)
	if err != nil {
		return nil, err
	}

	client.ctx = context.Background()
	k8sClient = &client
	return &client, nil
}
