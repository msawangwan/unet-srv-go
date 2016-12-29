package exception

import (
	"net/http"
)

type Catch func(http.ResponseWriter, *http.Request) *Handler

func (c Catch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := c(w, r); e != nil {
		http.Error(w, e.Message, e.Code)
	}
}
