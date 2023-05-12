package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v8"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type configuration struct {
	Secret string `env:"SECRET" envDefault:"xZp1CMxquSyj12h3TDyR"`
	Server string `env:"SERVER" envDefault:"https://localhost:8080"`
	Domain string `env:"DOMAIN" envDefault:"example.com"`
}

var Clientset *kubernetes.Clientset
var config *rest.Config
var Config configuration

func init() {

	err := InitConfig()
	if err != nil {
		log.Println(err)
	}
	err = ClusterConfig()
	if err != nil {
		log.Println(err)
	}
	log.Println("The controller is started ! ")

}

func InitConfig() error {
	err := env.Parse(&Config)
	if err != nil {
		return err
	}
	return nil

}

func ClusterConfig() error {
	var err error

	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")

	flag.Parse()

	if *kubeconfig != "" {
		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return err
		}

		// create the Clientset
		Clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			return err
		}
		// Retrieve the CA certificate data

	} else {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			return err
		}
		// creates the Clientset
		Clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			return err
		}

	}

	return nil
}
