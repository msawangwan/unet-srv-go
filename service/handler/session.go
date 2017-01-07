package handler

import (
	"net"

	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet/engine/session"
	"github.com/msawangwan/unet/env"
	"github.com/msawangwan/unet/service/exception"
)

// POST session/availability
func CheckSessionNameAvailable(g *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		la   *session.LobbyAvailability
		skey *session.Key
	)

	err := json.NewDecoder(r.Body).Decode(&skey)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	la = &session.LobbyAvailability{}
	err = la.CheckAvailability(skey.BareFormat, g.Pool, g.Log)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	json.NewEncoder(w).Encode(la)

	return nil
}

// POST session/new
func CreateNewSession(g *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		instance *session.Instance
		skey     *session.Key
	)

	err := json.NewDecoder(r.Body).Decode(&skey)

	instance, err = session.Create(skey.BareFormat, g.Pool, g.Log)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	json.NewEncoder(w).Encode(instance)

	return nil
}

// POST session/new/open
func MakeSessionActive(g *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		instance *session.Instance
		skey     *session.Key
	)

	err := json.NewDecoder(r.Body).Decode(&instance)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	key, err := instance.LoadSessionInstanceIntoMemory(g.Pool, g.Log)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	} else {
		if key != nil {
			skey = &session.Key{
				BareFormat:  instance.SessionID,
				RedisFormat: *key,
			}
		}
	}

	json.NewEncoder(w).Encode(skey)

	return nil
}

// POST session/new/join
func JoinExistingSession(g *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		instance *session.Instance
		skey     *session.Key
	)

	err := json.NewDecoder(r.Body).Decode(&skey)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	instance, err = session.Join(skey.BareFormat, g.Pool, g.Log)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	json.NewEncoder(w).Encode(instance)

	return nil
}

// POST session/new/instance/key
func KeyFromInstance(g *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		instance *session.Instance
		skey     *session.Key
	)

	err := json.NewDecoder(r.Body).Decode(&instance)
	if err != nil {
		return throw(err, err.Error(), 500)
	}

	key, err := instance.KeyFromInstance(g.Pool, g.Log)
	if err != nil {
		return throw(err, err.Error(), 500)
	} else {
		if key != nil {
			skey = &session.Key{
				BareFormat:  instance.SessionID,
				RedisFormat: *key,
			}
		}
	}

	json.NewEncoder(w).Encode(skey)

	return nil
}

// POST session/new/connect
func EstablishSessionConnection(g *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		owner *session.Owner
		conn  *session.Connection
	)

	err := json.NewDecoder(r.Body).Decode(&owner)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	var (
		ip string
	)

	ip = r.Header.Get("x-forwarded-for")
	if len(ip) == 0 { // we're proxying through nginx so we can prevent this
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return &exception.Handler{err, err.Error(), 500}
		}
	}

	result, key, err := session.EstablishConnection(owner.PlayerName, ip, g.Pool, g.Log)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	conn = &session.Connection{
		IsConnected: result,
		Address:     ip,
		Message:     *key,
	}

	json.NewEncoder(w).Encode(conn)

	return nil
}

// GET session/new/key
func RegisterNewSession(g *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		skey *int
	)

	skey, err := g.KeyGen.RegisterNewSession(g.Pool)
	if err != nil {
		return throw(err, err.Error(), 500)
	} else if skey == nil {
		return throw(nil, "nil key error", 500)
	} else {
		json.NewEncoder(w).Encode(
			struct {
				Value int `json:"value"`
			}{
				Value: *skey,
			},
		)
	}

	return nil
}

// GET session/active
func FetchAllActiveSessions(g *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		l *session.Lobby = &session.Lobby{}
	)

	err := l.PopulateLobbyList(g.Pool, g.Log)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	json.NewEncoder(w).Encode(l)

	return nil
}
