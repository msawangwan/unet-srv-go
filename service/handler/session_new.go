package handler

// service/handler/new_session.go handles session routes

import (
	//	"net"

	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet/engine/session"
	"github.com/msawangwan/unet/env"
	"github.com/msawangwan/unet/service/exception"
)

const (
	kLogPrefix_Session = "SESSION"
)

// RegisterNewSession creates a client session key, it is sent back as  json to
// the client and stored on the server in redis and in memory
//
// GET session/new/key
func RegisterNewSession(g *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		skey *int
		k    int
		ip   string
	)

	skey, err := g.KeyGen.RegisterNewSession(g.Pool)
	if err != nil {
		return throw(err, err.Error(), 500)
	} else if skey == nil {
		return throw(nil, "nil key error", 500)
	} else {
		k = *skey
		ip = r.Header.Get("x-forwarded-for")

		if len(ip) == 0 {
			ip = "invalid.client.add." + string(k) // TODO: handle for real
		}

		sh, _ := session.NewHandle(ip)
		if err != nil {
			g.Printf("tried to create a new handle but got an err: %s\n", err.Error())
		}
		g.SessionHandleManager.Add(k, sh)

		json.NewEncoder(w).Encode(
			struct {
				Value int `json:"value"`
			}{
				Value: k,
			},
		)
	}

	return nil
}

// HostNewGame starts a new game, given a request to host by a client
//
// POST session/host/instance
func HostNewGame(g *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	cleanup := setPrefix(kLogPrefix_Session, "HOST_NEW", g.Log)
	defer cleanup()

	j, err := parseJSONInt(r.Body)
	if err != nil {
		return throw(err, err.Error(), 500)
	} else if j == nil {
		g.Printf("nil key...") // TODO: handle
	}

	var (
		skey int    = *j
		ip   string = r.Header.Get("x-forwarded-for")

		shandle *session.Handle
	)

	shandle, err = g.SessionHandleManager.Get(skey)
	if err != nil {
		return throw(err, err.Error(), 500)
	}

	if ip == shandle.Owner {
		g.Printf("ip check ok") // TODO: do this sooner?
	}

	// NEED TO GET A GAME INSTANCE IN A SEPERATE ROUTE

	return nil
}
