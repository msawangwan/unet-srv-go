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
	j, err := parseJSON(r.Body)
	if err != nil {
		return raiseServerError(err)
	}

	defer g.PrefixReset()
	g.Prefix("handler", "client", "register")
	g.Printf("register new client handle")

	s, err := marshallJSONString(j)
	if err != nil {
		return raiseServerError(err)
	} else if s == nil {
		return raiseServerError(errors.New("nil string error"))
	}

	cname := *s
	//v, ok := j.(map[string]interface{})
	//if ok {
	//		cname, ok = v["value"].(string)
	//		if !ok {
	//			return raiseServerError(ErrFailedToRegisterClientHandle)
	//		}
	//	} else {
	//		return raiseServerError(ErrFailedToRegisterClientHandle)
	//	}

	k, err := g.KeyManager.GenerateNextClientID()
	if err != nil {
		return raiseServerError(err)
	}

	chid := *k

	err = session.RegisterClient(cname, chid, g.Pool, g.Log)
	if err != nil {
		return raiseServerError(err)
	}

	json.NewEncoder(w).Encode(
		struct {
			Value int `json:"value"`
		}{
			Value: chid,
		},
	)

	return nil
}

// RequestSessionKey : POST client/handle/host/key : return a session key for hosting
func RequestHostingKey(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	j, err := parseJSONInt(r.Body)
	if err != nil {
		return raiseServerError(err)
	} else if j == nil {
		return raiseServerError(errors.New("nil key in RequestSessionKey (line 70)"))
	}

	chid := *j

	defer g.PrefixReset()
	g.Prefix("handler", "client", "reqhostkey")
	g.Printf("session key has been requested by [clienthandle id: %d]", chid)

	var sessionHostKey int = -1

	mapped, err := session.IsMapped(chid, g.Pool, g.Log)
	if err != nil {
		return raiseServerError(err)
	} else if !mapped && err == nil {

		//g.Printf("client [handle id: %d]", *j)

		k, err := g.KeyManager.GenerateNextSessionKey()
		if err != nil {
			return raiseServerError(err)
		}

		sessionHostKey := *k

		err = session.MapToClient(chid, sessionHostKey, g.Pool, g.Log)
		if err != nil {
			return raiseServerError(err)
		}
	}

	json.NewEncoder(w).Encode(
		struct {
			Value int `json:"value"`
		}{
			Value: sessionHostKey,
		},
	)

	return nil
}

func RequestJoinKey(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	return nil
}
