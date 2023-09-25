package kubernetes

import (
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	istio "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kube "k8s.io/client-go/kubernetes"
	gatewayapiclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

type K8sClientInterface interface {
	Istio() istio.Interface
	GatewayAPI() gatewayapiclient.Interface
	ClientSet() kube.Interface

	GetMySQLPassword() (p string, err error)

	GetHttpFlowConfig(namespace string, host string) (*v1alpha3.DestinationRule, error)
	PutHttpFlowConfig(namespace string, host string, httpRule *networkingv1alpha3.ConnectionPoolSettings_HTTPSettings) (*v1alpha3.DestinationRule, error)
	DelHttpFlowConfig(namespace string, host string) error

	GetTcpFlowConfig(namespace string, host string) (*v1alpha3.DestinationRule, error)
	PutTcpFlowConfig(namespace string, host string, httpRule *networkingv1alpha3.ConnectionPoolSettings_TCPSettings) (*v1alpha3.DestinationRule, error)
	DelTcpFlowConfig(namespace string, host string) error

	GetDegradeConfig(namespace string, host string) (*v1alpha3.DestinationRule, error)
	PutDegradeConfig(namespace string, host string, degradeRule *networkingv1alpha3.OutlierDetection) (*v1alpha3.DestinationRule, error)
	DelDegradeConfig(namespace string, host string) error
}

func (in *K8SClient) ClientSet() kube.Interface {
	return in.k8s
}

func (in *K8SClient) Istio() istio.Interface {
	return in.istioClientset
}

func (in *K8SClient) GatewayAPI() gatewayapiclient.Interface {
	return in.gatewayapi
}

func (in *K8SClient) GetMySQLPassword() (string, error) {
	secret, err := in.k8s.CoreV1().Secrets("rbd-system").Get(in.ctx, "rbd-db", v1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(secret.Data["mysql-password"]), nil
}
