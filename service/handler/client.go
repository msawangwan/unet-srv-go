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

// RegisterClientHandle : POST client/handle/register : registers a new client handler
func RegisterClientHandle(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var (
		cname string
	)

	j, err := parseJSON(r.Body)
	if err != nil {
		return raiseServerError(err)
	}

	g.Prefix("handler", "client", "register")
	defer g.PrefixReset()

	g.Printf("register new client handle")

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

// GetSessionKey : POST client/handle/host/key : return a session key for hosting
func RequestHostingKey(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	g.Prefix("handler", "client", "reqhostkey")
	defer g.PrefixReset()

	g.Printf("new session key has been requested")

	j, err := parseJSONInt(r.Body)
	if err != nil {
		return raiseServerError(err)
	} else if j == nil {
		return raiseServerError(errors.New("nil key in GetSessionKey (line 72)"))
	}

	g.Printf("client [handle id: %d]", *j) // TODO: generate key rather than passing into func??

	sid, err := session.MapToClient(*j, g.SessionKeyGenerator, g.Pool, g.Log)
	if err != nil {
		return raiseServerError(err)
	}

	json.NewEncoder(w).Encode(
		struct {
			Value int `json:"value"`
		}{
			Value: *sid,
		},
	)

	return nil
}

func RequestJoinKey(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	return nil
}
