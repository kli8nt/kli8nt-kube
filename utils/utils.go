package utils

import (
	"context"
	"log"

	"github.com/kli8nt/kli8nt-kube/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func DeleteDeployment(namespace string, name string) error {
	err := config.Clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	log.Println("Deployment", name, "deleted from namespace", namespace)
	return nil
}

func GetDeployment(namespace, name string) (*appsv1.Deployment, error) {

	deploymentsClient := config.Clientset.AppsV1().Deployments(namespace)
	deployment, err := deploymentsClient.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return deployment, nil
}
func CreateDeployment(namespace string, deploymentName string, replicas int32, containerName string, image string) error {
	// Define the Deployment object
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentName,
				},
			},
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deploymentName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  containerName,
							Image: image,
						},
					},
				},
			},
		},
	}
	_, err := config.Clientset.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	log.Println("New Deployment :", deploymentName, "is created on namespace :", namespace)
	return nil
}

func UpdateDeployment(namespace string, deploymentName string, newImage string) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Get the Deployment object
		deployment, err := config.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if newImage != "" {
			deployment.Spec.Template.Spec.Containers[0].Image = newImage
		}

		// Update the Deployment object

		_, err = config.Clientset.AppsV1().Deployments(namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
		if err != nil {
			return err
		}

		log.Printf("Updated deployment %q in namespace %q.\n", deploymentName, namespace)

		return nil
	})
	return retryErr
}
