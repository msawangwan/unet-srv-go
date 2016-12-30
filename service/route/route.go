package route

import (
	"regexp"
	"sync"

	"net/http"

	"github.com/msawangwan/unet/env"
	"github.com/msawangwan/unet/service/handler"
	"github.com/msawangwan/unet/service/resource"
)

type route struct {
	Method  string
	Pattern *regexp.Regexp
	Handler http.Handler
}

type Table struct {
	Endpoints map[string]*route
	sync.Mutex
}

/*
routes:

	POST api/profile/availability -- is the profile name available for use?
	POST api/profile/create -- create a new profile
*/

func NewRouteTable(environment *env.Global) *Table {
	rt := &Table{
		Endpoints: map[string]*route{
			"availability": &route{
				Method:  "POST",
				Pattern: cache("api/profile/availability"),
				Handler: resource.Context{environment, handler.CheckProfileAvailability},
			},
			"profile_create": &route{
				Method:  "POST",
				Pattern: cache("api/profile/new"),
				Handler: resource.Context{environment, handler.CreateNewProfile},
			},
			"profile_world_new_data": &route{
				Method:  "POST",
				Pattern: cache("api/profile/world/load"),
				Handler: resource.Context{environment, handler.GenerateWorldData},
			},
		},
	}

	return rt
}

func cache(regex string) *regexp.Regexp {
	c, _ := regexp.Compile(regex)
	return c
}
