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
	j, err := parseJSON(r.Body)
	if err != nil {
		return raiseServerError(err)
	}

	s, err := marshallJSONString(j)
	if err != nil {
		return raiseServerError(err)
	}

	b, err := session.IsGameNameUnique(*s, g.Pool)
	if err != nil {
		return raiseServerError(err)
	}

	defer g.PrefixReset()
	g.Prefix("handler", "session", "verifyname")
	g.Printf("client requests to use [gamename: %s]", *s)
	g.Printf("gamename is unique: [%t]", b)

	json.NewEncoder(w).Encode(
		struct {
			Value bool `json:"value"`
		}{
			Value: b,
		},
	)

	return nil
}

// FetchLobby : GET session/handle/lobby/list
func FetchLobby(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	listing, err := session.GetLobby(g.Pool, g.Log)
	if err != nil {
		return raiseServerError(err)
	}

	g.Prefix("handler", "session", "fetchlobby")
	g.Printf("lobby listing: %q", listing)

	json.NewEncoder(w).Encode(
		struct {
			Listing []string `json:"listing"`
		}{
			Listing: listing,
		},
	)

	return nil
}

var (
	errJsonReadError = errors.New("[LoadGameHandler()] failed to parse client name from json (line 94)")
	errNilKey        = errors.New("[LoadGameHandler()] nil key (line: 68)")
)

// LoadGameHandler : POST session/handle/game/load/host
func LoadHostGameHandler(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	j, err := parseJSON(r.Body)
	if err != nil {
		return raiseServerError(err)
	}

	s, err := marshallJSONString(j)
	if err != nil {
		return raiseServerError(err)
	} else if s == nil {
		return raiseServerError(errJsonReadError)
	}

	var gname string = *s

	gk, err := g.KeyManager.GenerateNextGameID()
	if err != nil {
		return raiseServerError(err)
	} else if gk == nil {
		return raiseServerError(errNilKey)
	}

	var gid int = *gk

	gamehandlerstrp, err := game.CreateNewGame(gname, gid, g.Pool, g.Log)
	if err != nil {
		return raiseServerError(err)
	}

	err = g.Games.Add(*gamehandlerstrp)
	if err != nil {
		return raiseServerError(err)
	}

	g.Prefix("handler", "game", "loadhost")
	defer g.PrefixReset()

	g.Printf("recvd request to load host game")
	g.Printf("created a game key [key: %d]", gid)
	g.Printf("sending key to client ...")
	g.Printf("clients may now join this game using [key: %d]", gid)

	json.NewEncoder(w).Encode(
		struct {
			Value int `json:"value"`
		}{
			Value: gid,
		},
	)

	return nil
}

// LoadClientGameHandler : POST session/handle/game/load/client
func LoadClientGameHandler(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	j, err := parseJSON(r.Body)
	if err != nil {
		return raiseServerError(err)
	}

	s, err := marshallJSONString(j)
	if err != nil {
		return raiseServerError(err)
	}

	var gamename string = *s

	gameid, err := game.GetExistingGameByID(gamename, g.Pool, g.Log)
	if err != nil {
		return raiseServerError(err)
	}

	g.Prefix("handler", "game", "loadclient")
	defer g.PrefixReset()

	g.Printf("recvd request to load client game")
	g.Printf("client request to load [game name: %s] [game id: %d]", gamename, *gameid)

	json.NewEncoder(w).Encode(
		struct {
			Value int `json:"value"`
		}{
			Value: *gameid,
		},
	)

	return nil
}
