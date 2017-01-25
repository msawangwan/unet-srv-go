package gateway

import (
	"strings"

	"net/http"

	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/route"
)

// Multiplexer is the application gateway
type Multiplexer struct {
	*env.Global
	*route.Table
}

// NewMultiplexer is a factory that returns a new mux
func NewMultiplexer(environment *env.Global, routeTable *route.Table) *Multiplexer {
	if routeTable == nil {
		return &Multiplexer{
			Global: environment,
			Table:  route.NewRouteTable(environment),
		}
	} else {
		return &Multiplexer{
			Global: environment,
			Table:  routeTable,
		}
	}
}

// ServeHTTP implements http.HandlerFunc
func (mux *Multiplexer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		foundRoute bool   = false
		resource   string = r.URL.Path
	)

	// TODO: fix the logging here
	mux.Prefix("gateway", "multiplexer")
	defer mux.PrefixReset()

	ps := strings.Split(resource, "/")

	for path, route := range mux.Endpoints {
		if route.Pattern.MatchString(resource) == true {
			if route.Method == r.Method {
				rps := strings.Split(path, "/") // TODO: this might be too slow
				if rps[len(rps)-1] == ps[len(ps)-1] {
					foundRoute = true
					mux.Printf("found route: %s\n", resource)
					mux.Printf("serving resource at endpoint: %s\n", rps[len(rps)-1])
					route.Handler.ServeHTTP(w, r)
					break
				}
			}
		}
	}

	if !foundRoute {
		mux.Printf("invalid request: %s\n", resource)
	}
}
