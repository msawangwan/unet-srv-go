package resource

import (
	"net/http"

	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/exception"
)

// Context allows resource endpoints to access important environment variables
type Context struct {
	*env.Global
	Handle func(*env.Global, http.ResponseWriter, *http.Request) exception.Handler
}

// ServeHTTP implements http.Handler, should also handle all HTTP errors (todo)
func (c Context) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := c.Handle(c.Global, w, r); e != nil {
		c.SetPrefix("[RESOURCE][SERVE] ")
		defer c.SetPrefixDefault()
		c.Printf("got an error %s\n", e.Message)
		http.Error(w, e.Message, e.Code)
	}
}
