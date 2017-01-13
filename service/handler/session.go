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

// RegisterNewSession :  GET session/register/key
// DEPRECATE
func RegisterNewSession(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var (
		skey *int
		k    int
		ip   string
	)

	skey, err := g.SessionKeyGenerator.GenerateNextClientID()
	if err != nil {
		return raise(err, err.Error(), 500)
	} else if skey == nil {
		return raise(nil, "nil key error", 500)
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

// SetPlayerOwnerName : POST session/register/name
// DEPRECATE
func SetPlayerOwnerName(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	cleanup := setPrefix(logPrefixSession, "SET_NAME", g.Log)
	defer cleanup()

	j, err := parseJSON(r.Body)
	if err != nil {
		return raise(err, err.Error(), 500)
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
		return raise(err, err.Error(), 500)
	}

	shandle.SetOwner(name, r.Header.Get("x-forwarded-for"))

	return nil
}

// CheckGameNameAvailable : POST session/host/name/availability
func VerifyName(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	cleanup := setPrefix(logPrefixSession, "VERIFY_NAME", g.Log)
	defer cleanup()

	//var (
	//	la = &session.LobbyAvailability{}
	//)

	//j, err := parseJSON(r.Body)
	//if err != nil {
	//	return raiseServerError(err)
	//}

	// TODO: LEFT OF HERE

	//s, err := marshallJSONString(j)
	//if err != nil {
	//	return raiseServerError(err)
	//}

	g.Printf("CHECKING NAME AVAIL")

	//k := j.(map[string]interface{})["value"].(string)

	//if err = la.CheckAvailability(k, g.Pool, g.Log); err != nil {
	//	return raise(err, err.Error(), 500)
	//}

	//json.NewEncoder(w).Encode(la)

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
