package main

import (
	"github.com/msawangwan/unitywebservice/service"
	"github.com/msawangwan/unitywebservice/util"
	"log"
	"net/http"
)

func main() {
	util.Log.InitMessage("test server running ...")
	util.Log.InitMessage("listening on port 8080 ...")
	log.Fatal(http.ListenAndServe(":8080", service.ServiceGateway))
}
