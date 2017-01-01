package handler

import (
	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet/engine/session"
	"github.com/msawangwan/unet/env"
	"github.com/msawangwan/unet/model"
	"github.com/msawangwan/unet/service/exception"
)

func CheckSessionNameAvailable(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	var (
		q  *model.Key
		la *session.LobbyAvailability
	)

	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	la = &session.LobbyAvailability{}
	err = la.CheckAvailability(e, q.Query)
	if err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	json.NewEncoder(w).Encode(la)

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
