package service

import (
	"github.com/msawangwan/unitywebservice/util"
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
	var servedResource bool = false

	for _, route := range Route.Endpoints {
		if route.Pattern.MatchString(r.URL.Path) == true {
			if route.Method == r.Method {
				servedResource = true
				route.Handler(w, r)
				break
			}
		}
	}

	if !servedResource {
		util.Log.InvalidRequest(w, r)
	}
}
