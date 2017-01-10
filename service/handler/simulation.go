package handler

import (
	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet-srv-go/engine/game"
	"github.com/msawangwan/unet-srv-go/engine/session"
	"github.com/msawangwan/unet-srv-go/env"

	// "github.com/msawangwan/unet-srv-go/engine/game"
	// "github.com/msawangwan/unet-srv-go/engine/session"

	"github.com/msawangwan/unet-srv-go/service/exception"
)

// HostAndAttachNewSimulation : POST session/host/simulation
func HostAndAttachNewSimulation(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {
	cleanup := setPrefix(logPrefixSession, "HOST_NEW", g.Log)
	defer cleanup()

	j, err := parseJSON(r.Body)
	if err != nil {
		return raise(err, err.Error(), 500)
	}

	k := int(j.(map[string]interface{})["key"].(float64))
	label := j.(map[string]interface{})["value"].(string)

	var (
		shandle *session.Handle
	)

	shandle, err = g.SessionHandleManager.Get(k)
	if err != nil {
		return raise(err, err.Error(), 500)
	}

	if r.Header.Get("x-forwarded-for") == shandle.OwnerIP {
		g.Printf("ip check ok") // TODO: do this sooner?
	}

	var (
		sim *game.Simulation
	)

	sim, err = game.NewSimulation(label, game.GenerateSeedDebug(), g.GlobalError, g.Pool, g.Log)
	if err != nil {
		return raise(err, err.Error(), 500)
	}

	err = shandle.AttachSimulation(sim)
	if err != nil {
		return raise(err, err.Error(), 500)
	}

	json.NewEncoder(w).Encode(sim)

	return nil
}

// JoinGame : POST
func JoinGame(g *env.Global, w http.ResponseWriter, r *http.Request) exception.Handler {

	return nil
}
