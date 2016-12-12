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
	log.Printf("gateway attempt to match incoming request: %s\n", r.URL.Path)
	for _, route := range Routes {
		if route.Pattern.MatchString(r.URL.Path) == true {
			log.Printf("matched pattern!\n")
			if route.Method == r.Method {
				log.Printf("found a match, serving the request\n")
				route.Handler(w, r)
				return
			}
		}
	}
	log.Printf("failed to match on the requested resource\n")
	http.NotFound(w, r)
}
