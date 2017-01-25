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

// ServeHTTP implements http.Handler, should also handle all HTTP errors
func (c Context) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Prefix("context", "resource", "serveroute")

	defer func() {
		c.PrefixError("context", "resource", "fatal")
		c.SetLevelVerbose()

		if err := recover(); err != nil {
			c.Printf("caught a fatal error (panic)")
			c.Printf("%v", err)
		}

		c.PrefixReset()
		c.SetLevelDefault()
	}()

	c.Printf("calling handler mapped to: %s", r.URL.Path)

	if e := c.Handle(c.Global, w, r); e != nil {
		c.PrefixError("context", "resource", "error")
		c.Printf("%s", e.Print())
		c.PrefixReset()
	}
}
