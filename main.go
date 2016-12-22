package main

import (
	"github.com/msawangwan/unet/service"
	"github.com/msawangwan/unet/util"
	"log"
	"net/http"
)

func main() {
	util.Log.InitMessage("test server running ...")
	util.Log.InitMessage("listening on port 8080 ...")
	log.Fatal(http.ListenAndServe(":8080", service.ServiceGateway))
}
