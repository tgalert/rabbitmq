package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	//"log"
)

type User struct {
	name     string
	password string
}

func main() {
	host := flag.String("h", "localhost", "Domain name or IP address of RabbitMQ")
	flag.Parse()

	// Deploy RabbitMQ to Kubernetes cluster

	// Replace admin user
	baseUrl := "http://" + *host + ":15672/api"
	oldAdmin := User{"guest", "guest"}
	newAdmin := User{"admin", fmt.Sprintf("%s", uuid.New())}
	replaceAdmin(baseUrl, oldAdmin, newAdmin)

	// Save new admin user credentials in AWS Secrets Manager
}
