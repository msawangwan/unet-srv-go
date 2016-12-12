package main

import (
	"github.com/msawangwan/unitywebservice/service"
	"log"
	"net/http"
)

func main() {
	log.Printf("test server, serving routes from 8080...\n")
	//log.Fatal(http.ListenAndServe("8080", nil))
	log.Fatal(http.ListenAndServe(":8080", service.ServiceGateway))
}
