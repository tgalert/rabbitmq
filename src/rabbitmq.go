package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Create a new admin user, then delete the old admin user
func replaceAdmin(baseUrl string, oldAdmin User, newAdmin User) {

	log.Printf("Creating new admin: %v", newAdmin)
	makeRequest(
		http.MethodPut,
		baseUrl+"/users/"+newAdmin.name,
		fmt.Sprintf(`{"password":"%s","tags":"administrator"}`, newAdmin.password),
		oldAdmin,
	)

	log.Printf("Deleting old admin: %v", oldAdmin)
	makeRequest(
		http.MethodDelete,
		baseUrl+"/users/"+oldAdmin.name,
		"",
		oldAdmin,
	)
}

func makeRequest(method string, url string, body string, oldAdmin User) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		log.Fatal(err.Error())
	}
	req.SetBasicAuth(oldAdmin.name, oldAdmin.password)
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
