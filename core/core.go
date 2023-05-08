package core

import (
	"log"
	"time"

	"github.com/kli8nt/kli8nt-kube/config"
	"github.com/kli8nt/kli8nt-kube/utils"
)

func Start() {
	log.Println("Checking The connectivity with ", config.Config.Server)

	err := utils.CreateDeployment("default", "adam", 1, "nginx", "nginx")
	if err != nil {
		log.Println(err)
	}
	time.Sleep(10 * time.Second)

	err = utils.UpdateDeployment("default", "adam", "httpd")
	if err != nil {
		log.Println(err)
	}
	log.Println("The Deployment is Updated")

	time.Sleep(10 * time.Second)

	err = utils.DeleteDeployment("default", "adam")
	if err != nil {
		log.Println(err)
	}

}
