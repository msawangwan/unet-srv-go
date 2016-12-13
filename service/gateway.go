package service

import (
	"log"
	"net/http"
)

type gateway struct{}

var (
	ServiceGateway *gateway
)

func init() {
	ServiceGateway = &gateway{}
}

func (g *gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var foundResource bool = false

	for _, route := range Routes {
		if route.Pattern.MatchString(r.URL.Path) == true {
			if route.Method == r.Method {
				foundResource = true
				route.Handler(w, r)
				break
			}
		}
	}

	if foundResource {
		log.Printf("found a match, serving the request\n")
	} else {
		log.Printf("failed to match on the requested resource\n")
		http.NotFound(w, r)
	}
}
