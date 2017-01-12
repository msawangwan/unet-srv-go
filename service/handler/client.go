package handler

import (
	"errors"

	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet-srv-go/engine/session"
	"github.com/msawangwan/unet-srv-go/env"
	"github.com/msawangwan/unet-srv-go/service/exception"
)

const (
	logPrefixClient = "CLIENT"
)

var (
	ErrFailedToRegisterClientHandle = errors.New("failed to register client handle")
)

// RegisterClientHandle : POST client/handle/register : registers a new client
func RegisterClientHandle(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var (
		cname string
	)

	j, err := parseJSON(r.Body)
	if err != nil {
		return raiseServerError(err)
	}

	cleanup := setPrefix(logPrefixClient, "REGISTER_CLIENT_HANDLE", g.Log)
	defer cleanup()

	v, ok := j.(map[string]interface{})
	if ok {
		cname, ok = v["value"].(string)
		if !ok {
			return raiseServerError(ErrFailedToRegisterClientHandle)
		}
	} else {
		return raiseServerError(ErrFailedToRegisterClientHandle)
	}

	chid, err := session.RegisterClient(cname, g.SessionKeyGenerator, g.Pool, g.Log)
	if err != nil {
		return raiseServerError(err)
	}

	json.NewEncoder(w).Encode(
		struct {
			Value int `json:"value"`
		}{
			Value: *chid,
		},
	)

	return nil
}

func GetSessionKey(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	var (
		chID int
	)

	cleanup := setPrefix(logPrefixClient, "SESSION_KEY_REQ", g.Log)
	defer cleanup()

	g.Printf("new session key has been requested")

	j, err := parseJSONInt(r.Body)
	if err != nil {
		return raiseServerError(err)
	} else if j == nil {
		return raiseServerError(errors.New("nil key in GetSessionKey (line 72)"))
	}

	chID = *j
	// go into redis, and add a hash field for a session id

	json.NewEncoder(w).Encode(
		struct {
			Value int `json:"value"`
		}{
			Value: *j,
		},
	)

	return nil
}
