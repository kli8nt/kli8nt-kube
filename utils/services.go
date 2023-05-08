package utils

import (
	"context"
	"log"

	"github.com/kli8nt/kli8nt-kube/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

func ExposeDeployment(namespace string, deploymentName string) error {

	// Create the service object
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
			Labels: map[string]string{
				"app": deploymentName,
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
			Selector: map[string]string{
				"app": deploymentName,
			},
		},
	}

	// Create the service in the cluster
	_, err := config.Clientset.CoreV1().Services(namespace).Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	log.Printf("Service has been created for deployment with label app =%s in namespace %q.\n", deploymentName, namespace)

	return nil
}

func DeleteService(namespace, serviceName string) error {
	// Delete the service in the cluster
	err := config.Clientset.CoreV1().Services(namespace).Delete(context.Background(), serviceName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	log.Printf("Service %q has been deleted from namespace %q.\n", serviceName, namespace)

	return nil
}

func GetService(namespace, serviceName string) (*corev1.Service, error) {
	service, err := config.Clientset.CoreV1().Services(namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return service, nil
}
