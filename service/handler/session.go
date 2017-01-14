package handler

import (
	"errors"

	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/exception"

	"github.com/msawangwan/unet-srv-go/engine/game"
	"github.com/msawangwan/unet-srv-go/engine/session"
)

// VerifyName : POST session/handle/name/verification : check if unique name
func VerifyName(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	g.Prefix("handler", "session", "verifyname")
	defer g.PrefixReset()

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

	g.Prefix("handler", "session", "fetchlobby")
	g.Printf("lobby listing: %v", l)

	json.NewEncoder(w).Encode(l)

	return nil
}

// LoadGameHandler : POST session/handle/game
func LoadGameHandler(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var (
		id, gid int
		cid, gk *int
	)

	// client needs to send:
	// - client id (session id?? or can we get it from the cid??)

	cid, err := parseJSONInt(r.Body)
	if err != nil {
		return raiseServerError(err)
	} else if cid == nil {
		return raiseServerError(errors.New("nil key in LoadGameHandler (line 88)"))
	}

	id = *cid

	// server then:
	// enters a new instance in redis
	// returns a key to get the instance later

	gk, err = g.SessionKeyGenerator.GenerateNextGameID()
	if err != nil {
		return raiseServerError(err)
	} else if gk == nil {
		return raiseServerError(errors.New("[LoadGameHandler()] nil key (line: 68)"))
	}

	gid = *gk

	err = game.LoadNew(id, gid, g.Pool, g.Log)
	if err != nil {
		return raiseServerError(err)
	}

	g.Prefix("handler", "game", "loadgame")
	defer g.PrefixReset()

	g.Printf("recvd request to load game")
	g.Printf("created a game key [key: %d]", gid)
	g.Printf("sending key to client ...")

	json.NewEncoder(w).Encode(
		struct {
			Value int `json:"value"`
		}{
			Value: gid,
		},
	)

	// once loaded, a player can then join
	// only a host can load

	return nil
}
