package main

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

var clientset *kubernetes.Clientset

// Set up clients (must be called before the other functions)
func k8sInit(kubeconfig string) {
	log.Printf("Using kubeconfig file %s", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// Create service defined in YAML file (same as `kubectl apply -f file`)
func k8sCreateService(file string) {
	var spec coreV1.Service
	readFile(file, &spec)
	service, err := clientset.CoreV1().Services("default").Create(&spec)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created service %s from file %s", service.ObjectMeta.Name, file)
}

// Create deployment defined in YAML file (same as `kubectl apply -f file`)
func k8sCreateDeployment(file string) {
	var spec appsV1.Deployment
	readFile(file, &spec)
	deployment, err := clientset.AppsV1().Deployments("default").Create(&spec)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Created deployment %s from file %s", deployment.ObjectMeta.Name, file)
}

// Read the provided YAML file into the provided API object struct
func readFile(file string, obj interface{}) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = yaml.Unmarshal(bytes, obj)
	if err != nil {
		log.Fatal(err.Error())
	}
}
