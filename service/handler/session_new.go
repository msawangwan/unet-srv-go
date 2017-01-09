package handler

// service/handler/new_session.go handles session routes

import (
	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet/engine/game"
	"github.com/msawangwan/unet/engine/session"

	"github.com/msawangwan/unet/env"
	"github.com/msawangwan/unet/service/exception"
)

// - RegisterNewSession creates a client session key, it is sent back as json to
//     the client and stored on the server in redis and in memory
// - SetPlayerOwnerName assigns the players name to its associated session handle
// - CheckGameNameAvailable checks if the game name is valid (for host)
// - HostNewGame starts a new game, given a request to host by a client
// - FetchAllActive gets the lobby list

const (
	kLogPrefix_Session = "SESSION"
)

// GET session/register/key
func RegisterNewSession(g *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		skey *int
		k    int
		ip   string
	)

	skey, err := g.KeyGen.GenerateNext(g.Pool)
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

// POST session/register/name
func SetPlayerOwnerName(g *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	cleanup := setPrefix(kLogPrefix_Session, "SET_NAME", g.Log)
	defer cleanup()

	j, err := parseJSON(r.Body)
	if err != nil {
		return throw(err, err.Error(), 500)
	}

	var (
		k       int
		name    string
		shandle *session.Handle
	)

	k = int(j.(map[string]interface{})["key"].(float64)) // TODO: this will cause bugs, i just know it
	name = j.(map[string]interface{})["value"].(string)

	g.Printf("register sessionHandle [sessionID: %d] to owner: %s ...", k, name)
	g.Printf("success ...")

	shandle, err = g.SessionHandleManager.Get(k)
	if err != nil {
		return throw(err, err.Error(), 500)
	}

	shandle.SetOwner(name, r.Header.Get("x-forwarded-for"))

	return nil
}

// POST session/host/name/availability
func CheckGameNameAvailable(g *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	cleanup := setPrefix(kLogPrefix_Session, "CHECK_HOST_NAME", g.Log)
	defer cleanup()

	var (
		la *session.LobbyAvailability = &session.LobbyAvailability{}
	)

	j, err := parseJSON(r.Body)
	if err != nil {
		return throw(err, err.Error(), 500)
	}

	k := j.(map[string]interface{})["value"].(string)

	if err = la.CheckAvailability(k, g.Pool, g.Log); err != nil {
		return throw(err, err.Error(), 500)
	}

	json.NewEncoder(w).Encode(la)

	return nil
}

// POST session/host/instance
func HostNewGame(g *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	cleanup := setPrefix(kLogPrefix_Session, "HOST_NEW", g.Log)
	defer cleanup()

	j, err := parseJSON(r.Body)
	if err != nil {
		return throw(err, err.Error(), 500)
	}

	k := int(j.(map[string]interface{})["key"].(float64))
	label := j.(map[string]interface{})["value"].(string)

	var (
		shandle *session.Handle
	)

	shandle, err = g.SessionHandleManager.Get(k)
	if err != nil {
		return throw(err, err.Error(), 500)
	}

	if r.Header.Get("x-forwarded-for") == shandle.OwnerIP {
		g.Printf("ip check ok") // TODO: do this sooner?
	}

	var (
		sim *game.Simulation
	)

	sim, err = game.NewSimulation(label, game.GenerateSeedDebug(), g.Pool, g.Log)
	if err != nil {
		return throw(err, err.Error(), 500)
	}

	err = shandle.AttachSimulation(sim)
	if err != nil {
		return throw(err, err.Error(), 500)
	}

	json.NewEncoder(w).Encode(sim)

	return nil
}

// GET session/active
func FetchAllActiveSessions(g *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		l *session.Lobby = &session.Lobby{}
	)

	err := l.PopulateLobbyList(g.Pool, g.Log)
	if err != nil {
		return throw(err, err.Error(), 500)
	}

	g.Printf("lobby listing: %v", l)

	json.NewEncoder(w).Encode(l)

	return nil
}
