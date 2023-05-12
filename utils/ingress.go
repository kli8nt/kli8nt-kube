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

func DeleteTLSHost(ingressName, namespace, hostToDelete string) error {
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

	// If there is no TLS rule, return an error
	if tlsRule == nil {
		return err
	}

	// Remove the specified host from the hosts list of the first TLS rule
	newHosts := []string{}
	for _, host := range tlsRule.Hosts {
		if host != hostToDelete {
			newHosts = append(newHosts, host)
		}
	}
	tlsRule.Hosts = newHosts

	// Update the ingress object in the API server
	_, err = config.Clientset.NetworkingV1().Ingresses(namespace).Update(context.Background(), ingress, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	log.Println(hostToDelete, "is not in the TLS rules now!")

	return nil
}

func AddIngressRule(ingressName, namespace, serviceName string, port int32) error {

	// Get the ingress object from the API server
	ingress, err := config.Clientset.NetworkingV1().Ingresses(namespace).Get(context.Background(), ingressName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// create a new IngressRule object with the provided host name and backend service
	pathType := networkingv1.PathType("Prefix")
	newRule := networkingv1.IngressRule{
		Host: serviceName + "." + config.Config.Domain,
		IngressRuleValue: networkingv1.IngressRuleValue{
			HTTP: &networkingv1.HTTPIngressRuleValue{
				Paths: []networkingv1.HTTPIngressPath{
					{
						Path:     "/",
						PathType: &pathType, // specify the path type as Prefix
						Backend: networkingv1.IngressBackend{
							Service: &networkingv1.IngressServiceBackend{
								Name: serviceName, // specify the name of the backend service
								Port: networkingv1.ServiceBackendPort{
									Number: port, // specify the port number
								},
							},
						},
					},
				},
			},
		},
	}

	// append the new IngressRule object to the existing list of rules in the Ingress object
	ingress.Spec.Rules = append(ingress.Spec.Rules, newRule)

	// Update the ingress object in the API server
	_, err = config.Clientset.NetworkingV1().Ingresses(namespace).Update(context.Background(), ingress, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	log.Println(serviceName, "is exposed with ingress!")

	return nil
}

func DeleteIngressRule(ingressName, namespace, serviceName string) error {
	rules := []networkingv1.IngressRule{}

	// Get the ingress object from the API server
	ingress, err := config.Clientset.NetworkingV1().Ingresses(namespace).Get(context.Background(), ingressName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	for _, rule := range ingress.Spec.Rules {
		if rule.IngressRuleValue.HTTP.Paths[0].Backend.Service.Name != serviceName {
			rules = append(rules, rule)
		}
	}

	ingress.Spec.Rules = rules

	// Update the ingress object in the API server
	_, err = config.Clientset.NetworkingV1().Ingresses(namespace).Update(context.Background(), ingress, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	log.Println(serviceName, "is not exposed anymore in the ingress")

	return nil
}
