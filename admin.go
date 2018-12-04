package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	//"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type User struct {
	name     string
	password string
}

var adminOld User = User{"guest", "guest"}
var adminNew User = User{"admin", fmt.Sprintf("%s", uuid.New())}
var baseUrl string

func main() {
	host := flag.String("h", "localhost", "Public domain name or IP address of the RabbitMQ service")
	flag.Parse()
	baseUrl = "http://" + *host + ":15672/api"

	log.Printf("Creating new admin: %v", adminNew)
	createNewAdmin()
	log.Printf("Deleting old admin: %v", adminOld)
	deleteOldAdmin()
}

func createNewAdmin() {
	method := http.MethodPut
	url := baseUrl + "/users/" + adminNew.name
	body := fmt.Sprintf(`{"password":"%s","tags":"administrator"}`, adminNew.password)
	makeRequest(method, url, body)
}

func deleteOldAdmin() {
	method := http.MethodDelete
	url := baseUrl + "/users/" + adminOld.name
	makeRequest(method, url, "")
}

func makeRequest(method string, url string, body string) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		log.Fatal(err.Error())
	}
	req.SetBasicAuth(adminOld.name, adminOld.password)
	req.Header.Set("Content-Type", "application/json")
	log.Printf("Request: Method: %s, URL: %s, Body: %s", method, url, body)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Fatalf("Request failed: Status: %s, Response Body: %s", res.Status, body)
	}
	log.Print("Response: " + res.Status)
}
