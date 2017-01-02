package handler

import (
	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet/engine/session"
	"github.com/msawangwan/unet/env"
	"github.com/msawangwan/unet/service/exception"
)

func CheckSessionNameAvailable(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		info *session.Key
		la   *session.LobbyAvailability
	)

	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	la = &session.LobbyAvailability{}
	err = la.CheckAvailability(e, info.Info)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	json.NewEncoder(w).Encode(la)

	return nil
}

func CreateNewSession(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		info     *session.Key
		instance *session.Instance
	)

	err := json.NewDecoder(r.Body).Decode(&info)

	instance, err = session.Create(e, info.Info)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	json.NewEncoder(w).Encode(instance)

	return nil
}

func MakeSessionActive(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		instance *session.Instance
	)

	err := json.NewDecoder(r.Body).Decode(&instance)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	err = instance.LoadSessionInstanceIntoMemory(e)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	return nil
}

func JoinExistingSession(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		k        *session.Key
		instance *session.Instance
	)

	err := json.NewDecoder(r.Body).Decode(&k)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	instance, err = session.Join(e, k.Info)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	json.NewEncoder(w).Encode(instance)

	return nil
}

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
