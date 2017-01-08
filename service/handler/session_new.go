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
	var (
		skey int
	)

	err := json.NewDecoder(r.Body).Decode(skey)
	if err != nil {
		return throw(err, err.Error(), 500)
	}

	g.Printf("DECODED JSON INTO INT: %d", skey)

	return nil
}
