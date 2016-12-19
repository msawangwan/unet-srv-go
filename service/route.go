package service

import (
	"github.com/msawangwan/unitywebservice/util"
	"net/http"
	"regexp"
	"sync"
)

type route struct {
	Method  string
	Pattern *regexp.Regexp
	Handler http.HandlerFunc
}

type routeTable struct {
	Endpoints map[string]*route
	sync.Mutex
}

var (
	Route *routeTable
)

/*
routes:

	POST api/availability -- is the profile name available for use?
	POST api/profile/create -- create a new profile
*/

func init() {
	util.Log.InitMessage("compiling routes ...")

	Route = &routeTable{
		Endpoints: map[string]*route{
			"availability": &route{
				Method:  "POST",
				Pattern: Cache("api/availability"),
				Handler: util.Log.ResourceRequest(availability),
			},
			"profile_create": &route{
				Method:  "POST",
				Pattern: Cache("api/profile/create"),
				Handler: util.Log.ResourceRequest(profileCreate),
			},
			"profile_world_new_data": &route{
				Method:  "POST",
				Pattern: Cache("api/profile/world/new_data"),
				Handler: util.Log.ResourceRequest(profileNewWorldData),
			},
		},
	}

	util.Log.InitMessage("routes ready!")
}

func Cache(regex string) *regexp.Regexp {
	cache, _ := regexp.Compile(regex)
	return cache
}
