package handler

// service/handler/new_session.go handles session routes

import (
	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet-srv-go/engine/session"

	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/exception"
)

const (
	logPrefixSession = "SESSION"
)

// VerifyName : POST session/handle/name/verification : check if unique name
func VerifyName(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	cleanup := setPrefix(logPrefixSession, "VERIFY_NAME", g.Log)
	defer cleanup()

	j, err := parseJSON(r.Body)
	if err != nil {
		return raiseServerError(err)
	}

	s, err := marshallJSONString(j)
	if err != nil {
		return raiseServerError(err)
	}

	g.Printf("client requests to use [gamename: %s]", *s)

	b, err := session.IsGameNameUnique(*s, g.Pool)
	if err != nil {
		return raiseServerError(err)
	}

	g.Printf("gamename is unique: [%t]", b)

	if b {
		err = session.PostGameToLobby(*s, g.Pool)
		if err != nil {
			return raiseServerError(err)
		} else {
			g.Printf("posted session to lobby")
		}
	}

	json.NewEncoder(w).Encode(
		struct {
			Value bool `json:"value"`
		}{
			Value: b,
		},
	)

	return nil
}

// FetchAllActiveSessions : GET session/join/lobby/list
func FetchAllActiveSessions(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var (
		l = &session.Lobby{}
	)

	err := l.PopulateLobbyList(g.Pool, g.Log)
	if err != nil {
		return raise(err, err.Error(), 500)
	}

	g.Printf("lobby listing: %v", l)

	json.NewEncoder(w).Encode(l)

	return nil
}
