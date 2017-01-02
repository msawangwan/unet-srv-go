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
func CheckSessionNameAvailable(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		la   *session.LobbyAvailability
		skey *session.SessionKey
	)

	err := json.NewDecoder(r.Body).Decode(&skey)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	la = &session.LobbyAvailability{}
	err = la.CheckAvailability(e, skey.BareFormat)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	json.NewEncoder(w).Encode(la)

	return nil
}

// POST session/new
func CreateNewSession(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		instance *session.Instance
		skey     *session.SessionKey
	)

	err := json.NewDecoder(r.Body).Decode(&skey)

	instance, err = session.Create(e, skey.BareFormat)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	json.NewEncoder(w).Encode(instance)

	return nil
}

// POST session/new/open
func MakeSessionActive(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		instance *session.Instance
		skey     *session.SessionKey
	)

	err := json.NewDecoder(r.Body).Decode(&instance)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	key, err := instance.LoadSessionInstanceIntoMemory(e)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	} else {
		if key != nil {
			skey = &session.SessionKey{
				BareFormat:  instance.SessionID,
				RedisFormat: *key,
			}
		}
	}

	json.NewEncoder(w).Encode(skey)

	return nil
}

// POST session/new/join
func JoinExistingSession(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		instance *session.Instance
		skey     *session.SessionKey
	)

	err := json.NewDecoder(r.Body).Decode(&skey)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	instance, err = session.Join(e, skey.BareFormat)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	json.NewEncoder(w).Encode(instance)

	return nil
}

// POST session/new/connect
func EstablishSessionConnection(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		instance *session.Instance
		conn     *session.Connection
	)

	err := json.NewDecoder(r.Body).Decode(&instance)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	result, key, err := instance.Connect(e, ip)
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

// GET session/active
func FetchAllActiveSessions(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		l *session.Lobby = &session.Lobby{}
	)

	err := l.PopulateLobbyList(e)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	json.NewEncoder(w).Encode(l)

	return nil
}
