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
	c.SetPrefix("[RESOURCE][SERVE] ")

	defer func() {
		c.SetPrefix("[RESOURCE][ERROR][CAUGHT_FATAL] ")
		c.SetLevelVerbose()

		if err := recover(); err != nil {
			c.Printf("caught a fatal error (panic)")
			c.Printf("%v", err)
		}

		c.SetPrefixDefault()
		c.SetLevelDefault()
	}()

	c.Printf("calling handler function mapped to: %s", r.URL.Path)

	if e := c.Handle(c.Global, w, r); e != nil {
		c.SetPrefix("[RESOURCE][ERROR] ")
		c.Printf("%s", e.Print())
		c.SetPrefixDefault()
	}
}
