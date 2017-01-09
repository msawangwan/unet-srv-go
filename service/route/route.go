package route

import (
	"regexp"
	"sync"

	"net/http"

	"github.com/msawangwan/unet/env"
	"github.com/msawangwan/unet/service/handler"
	"github.com/msawangwan/unet/service/resource"
)

const (
	kPrefixFallback = "api" // TODO: get from config
)

// type route consists of:
// a method
// a regex pattern (path)
// a handler
type route struct {
	Method  string
	Pattern *regexp.Regexp
	Handler http.Handler
}

// type Table represents a map of routes
type Table struct {
	Endpoints map[string]*route
	sync.Mutex
}

// NewRouteTable creates a table of all available routes, addes them to a regex
// cache and assigns their respective handlers
func NewRouteTable(environment *env.Global) *Table {
	rt := &Table{
		Endpoints: map[string]*route{
			"profile/availability": &route{
				Method:  "POST",
				Pattern: cache("profile/availability"),
				Handler: resource.Context{environment, handler.CheckProfileAvailability},
			},
			"profile/new": &route{
				Method:  "POST",
				Pattern: cache("profile/new"),
				Handler: resource.Context{environment, handler.CreateNewProfile},
			},
			"profile/world/load": &route{
				Method:  "POST",
				Pattern: cache("profile/world/load"),
				Handler: resource.Context{environment, handler.GenerateWorldData},
			},
			"session/active": &route{
				Method:  "GET",
				Pattern: cache("session/active"),
				Handler: resource.Context{environment, handler.FetchAllActiveSessions},
			},
			"session/availability": &route{
				Method:  "POST",
				Pattern: cache("session/availability"),
				Handler: resource.Context{environment, handler.CheckSessionNameAvailable},
			},
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
			"session/host/instance": &route{
				Method:  "POST",
				Pattern: cache("session/host/instance"),
				Handler: resource.Context{environment, handler.HostNewGame},
			},
			"session/new": &route{
				Method:  "POST",
				Pattern: cache("session/new"),
				Handler: resource.Context{environment, handler.CreateNewSession},
			},
			"session/new/open": &route{
				Method:  "POST",
				Pattern: cache("session/new/open"),
				Handler: resource.Context{environment, handler.MakeSessionActive},
			},
			"session/new/join": &route{
				Method:  "POST",
				Pattern: cache("session/new/join"),
				Handler: resource.Context{environment, handler.JoinExistingSession},
			},
			"session/new/connect": &route{
				Method:  "POST",
				Pattern: cache("session/new/connect"),
				Handler: resource.Context{environment, handler.EstablishSessionConnection},
			},
			"session/new/instance/key": &route{
				Method:  "POST",
				Pattern: cache("session/new/instance/key"),
				Handler: resource.Context{environment, handler.KeyFromInstance},
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
	c, _ := regexp.Compile(kPrefixFallback + "/" + regex)
	return c
}
