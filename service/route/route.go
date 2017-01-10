package route

import (
	"regexp"
	"sync"

	"net/http"

	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/handler"
	"github.com/msawangwan/unet-srv-go/service/resource"
)

const (
	prefixFallback = "api" // TODO: get from config
)

// route is defined by an HTTP method, regex pattern (url) and a handler function
type route struct {
	Method  string
	Pattern *regexp.Regexp
	Handler http.Handler
}

// Table represents a map of routes
type Table struct {
	Endpoints map[string]*route
	sync.Mutex
}

// NewRouteTable creates a table of all available routes, addes them to a regex
// cache and assigns their respective handlers
func NewRouteTable(environment *env.Global) *Table {
	rt := &Table{
		Endpoints: map[string]*route{
			"session/register/key": &route{
				Method:  "GET",
				Pattern: cache("session/register/key"),
				Handler: resource.Context{environment, handler.RegisterNewSession},
			},
			"session/register/name": &route{
				Method:  "POST",
				Pattern: cache("session/register/name"),
				Handler: resource.Context{environment, handler.SetPlayerOwnerName},
			},
			"session/host/name/availability": &route{
				Method:  "POST",
				Pattern: cache("session/host/name/availability"),
				Handler: resource.Context{environment, handler.CheckGameNameAvailable},
			},
			"session/host/simulation": &route{
				Method:  "POST",
				Pattern: cache("session/host/simulation"),
				Handler: resource.Context{environment, handler.HostAndAttachNewSimulation},
			},
			"session/lobby/list": &route{
				Method:  "GET",
				Pattern: cache("session/active/list"),
				Handler: resource.Context{environment, handler.FetchAllActiveSessions},
			},
			"game/update/start": &route{
				Method:  "POST",
				Pattern: cache("game/update/start"),
				Handler: resource.Context{environment, handler.StartGameUpdate},
			},
			"game/update/enter": &route{
				Method:  "POST",
				Pattern: cache("game/update/enter"),
				Handler: resource.Context{environment, handler.EnterGameUpdate},
			},
			"game/update/frame": &route{
				Method:  "POST",
				Pattern: cache("game/update/frame"),
				Handler: resource.Context{environment, handler.GameFrameUpdate},
			},
		},
	}

	return rt
}

func cache(regex string) *regexp.Regexp {
	c, _ := regexp.Compile(prefixFallback + "/" + regex)
	return c
}
