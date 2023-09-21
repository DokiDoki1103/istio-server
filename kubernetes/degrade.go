package kubernetes

import (
	"fmt"
	"istio-server/config"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (in *K8SClient) GetDegradeConfig(namespace string, host string) (*v1alpha3.DestinationRule, error) {
	return in.Istio().
		NetworkingV1alpha3().
		DestinationRules(namespace).
		Get(in.ctx, config.Degrade+"-"+host, v1.GetOptions{})
}

func (in *K8SClient) PutDegradeConfig(namespace string, host string, degradeRule *networkingv1alpha3.OutlierDetection) (*v1alpha3.DestinationRule, error) {
	name := config.Degrade + "-" + host
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
				"app":  config.Degrade,
			},
		},
		Spec: networkingv1alpha3.DestinationRule{
			Host: fmt.Sprintf("%s.%s.svc.cluster.local", host, namespace),
			TrafficPolicy: &networkingv1alpha3.TrafficPolicy{
				OutlierDetection: degradeRule,
			},
		},
	}
	tx := in.Istio().NetworkingV1alpha3().DestinationRules(namespace)

	get, err := tx.Get(in.ctx, name, v1.GetOptions{})
	if err != nil {
		return tx.Create(in.ctx, r, v1.CreateOptions{})
	}
	get.Spec.TrafficPolicy.OutlierDetection = degradeRule
	return tx.Update(in.ctx, get, v1.UpdateOptions{})
}

func (in *K8SClient) DelDegradeConfig(namespace string, name string) error {
	return in.Istio().
		NetworkingV1alpha3().
		DestinationRules(namespace).
		Delete(in.ctx, config.Degrade+"-"+name, v1.DeleteOptions{})
}
