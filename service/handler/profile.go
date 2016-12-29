package handler

import (
	"encoding/json"
	"net/http"

	"github.com/msawangwan/unet/env"
	"github.com/msawangwan/unet/model"
	"github.com/msawangwan/unet/service/exception"
)

// Availability expects a struct containing a requested profile name and returns
// false if the name is already in use
func CheckProfileAvailability(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	if r.Body == nil {
		return &exception.Handler{errNilBody, "couldn't check for a name, did not recieve a body", 500}
	}

	var ps model.ProfileSearch

	err := json.NewDecoder(r.Body).Decode(&ps)
	if err != nil {
		return &exception.Handler{err, "error decoding json", 500}
	}

	_, err = ps.IsProfileNameAvailable(e.RedisHandle)
	if err != nil {
		return &exception.Handler{err, "error querying for profile name availability", 500}
	}

	json.NewEncoder(w).Encode(&ps)

	return nil
}

// ProfileCreate will generate a new profile given a valid name, create a new
// entry in the database as well as make an entry in the redis store
func CreateNewProfile(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	if r.Body == nil {
		return &exception.Handler{errNilBody, "got a nil body", 500}
	}

	var gamestate *model.GameState
	var starmap *model.StarMap
	var profile *model.Profile
	var n model.ProfileName

	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		return &exception.Handler{err, "error decoding json", 500}
	}

	profile, err = model.CreateNewProfile(n.Text, e.PostgreHandle)
	if err != nil {
		return &exception.Handler{err, "error creating a new profile in the db", 500}
	}

	err = profile.MarkNameAsNotAvailable(e.RedisHandle)
	if err != nil {
		return &exception.Handler{err, "error marking name as unavailable", 500}
	}

	starmap = model.NewMapDefaultParams(profile.Seed)
	gamestate = model.NewGameState(profile, starmap)

	json.NewEncoder(w).Encode(gamestate)

	return nil
}

// GenerateWorldData creates new world
func GenerateWorldData(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	if r.Body == nil {
		return &exception.Handler{errNilBody, "got a nil body", 500}
	}

	var sd model.StarData

	err := json.NewDecoder(r.Body).Decode(&sd)
	if err != nil {
		return &exception.Handler{err, "error decoding json", 500}
	}

	json.NewEncoder(w).Encode(&sd)

	return nil
}
