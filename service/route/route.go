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
func NewRouteTable(globals *env.Global) *Table {
	rt := &Table{
		Endpoints: map[string]*route{
			"client/handle/register": &route{
				Method:  "POST",
				Pattern: cache("client/handle/register"),
				Handler: resource.Context{globals, handler.RegisterClientHandle},
			},
			"client/handle/session/key": &route{
				Method:  "POST",
				Pattern: cache("client/handle/session/key"),
				Handler: resource.Context{globals, handler.RequestHostingKey},
			},
			"session/handle/name/verification": &route{
				Method:  "POST",
				Pattern: cache("session/handle/name/verification"),
				Handler: resource.Context{globals, handler.VerifyName},
			},
			"session/handle/lobby/fetch": &route{
				Method:  "",
				Pattern: cache(""),
				Handler: resource.Context{globals, handler.FetchLobby},
			},
			"session/handle/game/load": &route{
				Method:  "POST",
				Pattern: cache("session/handle/game/load"),
				Handler: resource.Context{globals, handler.LoadGameHandler},
			},
			"game/world/load": &route{
				Method:  "POST",
				Pattern: cache("game/world/load"),
				Handler: resource.Context{globals, handler.LoadWorld},
			},
			"game/world/join": &route{
				Method:  "POST",
				Pattern: cache("game/world/join"),
				Handler: resource.Context{globals, handler.JoinGameWorld},
			},
			"poll/start": &route{
				Method:  "POST",
				Pattern: cache("poll/start"),
				Handler: resource.Context{globals, handler.PollStart},
			},
			"poll/update": &route{
				Method:  "POST",
				Pattern: cache("poll/update"),
				Handler: resource.Context{globals, handler.PollUpdate},
			},
		},
	}

	return rt
}

func cache(regex string) *regexp.Regexp {
	c, _ := regexp.Compile(prefixFallback + "/" + regex)
	return c
}
