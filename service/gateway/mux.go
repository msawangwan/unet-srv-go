package gateway

import (
	"net/http"

	"github.com/msawangwan/unet/env"
	"github.com/msawangwan/unet/service/route"
)

type Multiplexer struct {
	*env.Global
	*route.Table
}

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

func (mux *Multiplexer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var foundRoute bool = false

	for _, route := range mux.Endpoints {
		if route.Pattern.MatchString(r.URL.Path) == true {
			if route.Method == r.Method {
				foundRoute = true
				mux.Printf("found route: %s\n", r.URL.Path)
				route.Handler.ServeHTTP(w, r)
				break
			}
		}
	}

	if !foundRoute {
		mux.Printf("invalid request: %s\n", r.URL.Path)
	}
}
