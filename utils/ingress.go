package utils

import (
	"context"
	"log"

	"github.com/kli8nt/kli8nt-kube/config"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AddTLSHost(ingressName, namespace, newHost string) error {
	// Get the ingress object from the API server
	ingress, err := config.Clientset.NetworkingV1().Ingresses(namespace).Get(context.Background(), ingressName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Get the first TLS rule (if any) in the ingress object
	var tlsRule *networkingv1.IngressTLS
	if len(ingress.Spec.TLS) > 0 {
		tlsRule = &ingress.Spec.TLS[0]
	}

	// If there is no TLS rule yet, create a new one
	if tlsRule == nil {
		tlsRule = &networkingv1.IngressTLS{}
		ingress.Spec.TLS = append(ingress.Spec.TLS, *tlsRule)
	}

	// Append the new host to the hosts list of the first TLS rule
	tlsRule.Hosts = append(tlsRule.Hosts, newHost)

	// Update the ingress object in the API server
	_, err = config.Clientset.NetworkingV1().Ingresses(namespace).Update(context.Background(), ingress, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	log.Println(newHost, "is added is secured now")

	return nil
}
