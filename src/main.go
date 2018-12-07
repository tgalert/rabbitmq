package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	//"log"
	"os"
)

type User struct {
	name     string
	password string
}

func main() {
	kubeconfig := flag.String("k", os.Getenv("HOME")+"/.kube/config", "kubeconfig file")
	flag.Parse()

	// Deploy RabbitMQ to Kubernetes cluster
	k8sInit(*kubeconfig)
	k8sCreateService("k8s/service.yml")
	k8sCreateDeployment("k8s/deployment.yml")
	os.Exit(0)

	host := "localhost"
	// Replace admin user
	baseUrl := "http://" + host + ":15672/api"
	oldAdmin := User{"guest", "guest"}
	newAdmin := User{"admin", fmt.Sprintf("%s", uuid.New())}
	replaceAdmin(baseUrl, oldAdmin, newAdmin)

	// Save new admin user credentials in AWS Secrets Manager
}
