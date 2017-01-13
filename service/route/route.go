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
			"client/handle/host/key": &route{
				Method:  "POST",
				Pattern: cache("client/handle/host/key"),
				Handler: resource.Context{globals, handler.RequestHostingKey},
			},
			"client/handle/join/key": &route{
				Method:  "POST",
				Pattern: cache("client/handle/join/key"),
				Handler: resource.Context{globals, handler.SetPlayerOwnerName},
			},
			"session/handle/name/verification": &route{
				Method:  "POST",
				Pattern: cache("session/handle/name/verification"),
				Handler: resource.Context{globals, handler.VerifyName},
			},
		},
	}

	return rt
}

func cache(regex string) *regexp.Regexp {
	c, _ := regexp.Compile(prefixFallback + "/" + regex)
	return c
}
