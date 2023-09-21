package kubernetes

import (
	"fmt"
	"istio-server/config"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (in *K8SClient) GetTcpFlowConfig(namespace string, host string) (*v1alpha3.DestinationRule, error) {
	return in.Istio().
		NetworkingV1alpha3().
		DestinationRules(namespace).
		Get(in.ctx, config.TcpFlow+"-"+host, v1.GetOptions{})
}

func (in *K8SClient) PutTcpFlowConfig(namespace string, host string, tcpRule *networkingv1alpha3.ConnectionPoolSettings_TCPSettings) (*v1alpha3.DestinationRule, error) {
	name := config.TcpFlow + "-" + host
	r := &v1alpha3.DestinationRule{
		TypeMeta: v1.TypeMeta{
			Kind:       config.DestinationRule,
			APIVersion: config.V1alpha3ApiVersion,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"name": host,
				"app":  config.TcpFlow,
			},
		},
		Spec: networkingv1alpha3.DestinationRule{
			Host: fmt.Sprintf("%s.%s.svc.cluster.local", host, namespace),
			TrafficPolicy: &networkingv1alpha3.TrafficPolicy{
				ConnectionPool: &networkingv1alpha3.ConnectionPoolSettings{
					Tcp: tcpRule,
				},
			},
		},
	}
	tx := in.Istio().NetworkingV1alpha3().DestinationRules(namespace)

	get, err := tx.Get(in.ctx, name, v1.GetOptions{})
	if err != nil {
		return tx.Create(in.ctx, r, v1.CreateOptions{})
	}
	get.Spec.TrafficPolicy.ConnectionPool.Tcp = tcpRule
	return tx.Update(in.ctx, get, v1.UpdateOptions{})

}
func (in *K8SClient) DelTcpFlowConfig(namespace string, host string) error {
	return in.Istio().
		NetworkingV1alpha3().
		DestinationRules(namespace).
		Delete(in.ctx, config.TcpFlow+"-"+host, v1.DeleteOptions{})
}

func (in *K8SClient) GetHttpFlowConfig(namespace string, host string) (*v1alpha3.DestinationRule, error) {

	return in.Istio().
		NetworkingV1alpha3().
		DestinationRules(namespace).
		Get(in.ctx, config.HttpFlow+"-"+host, v1.GetOptions{})
}

func (in *K8SClient) PutHttpFlowConfig(namespace string, host string, httpRule *networkingv1alpha3.ConnectionPoolSettings_HTTPSettings) (*v1alpha3.DestinationRule, error) {
	name := config.HttpFlow + "-" + host
	r := &v1alpha3.DestinationRule{
		TypeMeta: v1.TypeMeta{
			Kind:       config.DestinationRule,
			APIVersion: config.V1alpha3ApiVersion,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"name": host,
				"app":  config.HttpFlow,
			},
		},
		Spec: networkingv1alpha3.DestinationRule{
			Host: fmt.Sprintf("%s.%s.svc.cluster.local", host, namespace),
			TrafficPolicy: &networkingv1alpha3.TrafficPolicy{
				ConnectionPool: &networkingv1alpha3.ConnectionPoolSettings{
					Http: httpRule,
				},
			},
		},
	}
	tx := in.Istio().NetworkingV1alpha3().DestinationRules(namespace)

	get, err := tx.Get(in.ctx, name, v1.GetOptions{})
	if err != nil {
		return tx.Create(in.ctx, r, v1.CreateOptions{})
	}
	get.Spec.TrafficPolicy.ConnectionPool.Http = httpRule
	return tx.Update(in.ctx, get, v1.UpdateOptions{})
}

func (in *K8SClient) DelHttpFlowConfig(namespace string, host string) error {
	return in.Istio().
		NetworkingV1alpha3().
		DestinationRules(namespace).
		Delete(in.ctx, config.HttpFlow+"-"+host, v1.DeleteOptions{})
}
