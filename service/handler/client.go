package handler

import (
	"errors"

	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet-srv-go/engine/session"
	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/exception"
)

var (
	ErrFailedToRegisterClientHandle = errors.New("failed to register client handle")
)

// RegisterClientHandle : POST client/handle/register : registers a new client
func RegisterClientHandle(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var (
		cname string
	)

	j, err := parseJSON(r.Body)
	if err != nil {
		return raiseServerError(err)
	}

	cleanup := setPrefix("CLIENT", "REGISTER_CLIENT_HANDLE", g.Log)
	defer cleanup()

	g.Printf("registered client ...")

	g.Printf("%v", j)

	v, ok := j.(map[string]interface{})
	if ok {
		cname, ok = v["value"].(string)
		if !ok {
			return raiseServerError(ErrFailedToRegisterClientHandle)
		}
	} else {
		return raiseServerError(ErrFailedToRegisterClientHandle)
	}

	chid, err := session.RegisterClient(cname, g.SessionKeyGenerator, g.Pool, g.Log)
	if err != nil {
		return raiseServerError(err)
	}

	json.NewEncoder(w).Encode(
		struct {
			Value int `json:"value"`
		}{
			Value: *chid,
		},
	)

	return nil
}
