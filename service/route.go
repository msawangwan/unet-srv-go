package service

import (
	"log"
	"net/http"
	"regexp"
)

type route struct {
	Method  string
	Pattern *regexp.Regexp
	Handler http.HandlerFunc
}

type routeTable map[string]*route

var (
	Routes routeTable
)

/*
routes:

	POST api/availability -- is the profile name available for use?
*/

func init() {
	log.Printf("init route table ...\n")

	Routes = map[string]*route{
		"availability": &route{
			Method:  "POST",
			Pattern: Cache("api/availability"),
			Handler: Log.resourceRequest(availability),
		},
	}

	log.Printf("route table init success ...\n")
}

func Cache(regex string) *regexp.Regexp {
	cache, _ := regexp.Compile(regex)
	return cache
}
