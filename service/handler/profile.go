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
	if err := checkBody(r); err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	var ps model.ProfileSearch

	err := json.NewDecoder(r.Body).Decode(&ps)
	if err != nil {
		return &exception.Handler{err, "checkProfileAvailability error", 500}
	} else {
		e.Printf("checking profile name availability: %s\n", ps.Name)
	}

	_, err = ps.IsProfileNameAvailable(e.RedisHandle)
	if err != nil {
		return &exception.Handler{err, "error querying for profile name availability", 500}
	} else {
		e.Printf("name availability: %+v", ps.IsAvailable)
	}

	json.NewEncoder(w).Encode(&ps)

	return nil
}

// ProfileCreate will generate a new profile given a valid name, create a new
// entry in the database as well as make an entry in the redis store
func CreateNewProfile(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	if err := checkBody(r); err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	var (
		gamestate *model.GameState
		starmap   *model.StarMap
		profile   *model.Profile
		n         model.ProfileName
	)

	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		return &exception.Handler{err, "error decoding json", 500}
	}

	profile, err = model.CreateNewProfile(n.Text, e.PostgreHandle)
	if err != nil {
		return &exception.Handler{err, "error creating a new profile in the db", 500}
	} else {
		e.Printf("created a new profile:\n")
		e.Printf("\t%+v\n", profile)
	}

	err = profile.MarkNameAsNotAvailable(e.RedisHandle)
	if err != nil {
		return &exception.Handler{err, "error marking name as unavailable", 500}
	} else {
		e.Printf("marked name in use:\n")
		e.Printf("\t%s\n", profile.Name)
	}

	starmap = model.NewMapDefaultParams(profile.Seed)
	gamestate = model.NewGameState(profile, starmap)

	e.Printf("new starmap:\n")
	e.Printf("\t%+v\n", starmap)
	e.Printf("new gamestate:\n")
	e.Printf("\t%+v\n", gamestate)

	json.NewEncoder(w).Encode(gamestate)

	return nil
}

// GenerateWorldData creates new world, expects json "profile"
func GenerateWorldData(e *env.Global, w http.ResponseWriter, r *http.Request) *exception.Handler {
	if err := checkBody(r); err != nil {
		return &exception.Handler{err, err.Error(), 500}
	}

	var (
		p *model.Profile
	)

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		return &exception.Handler{err, "error decoding json", 500}
	} else {
		e.Printf("got request to load profile into memory ...\n")
		e.Printf("\t%+v\n", p)
	}

	if err = p.LoadIntoMemory(e); err != nil {
		return &exception.Handler{err, "error loading the profile world into memory", 500}
	} else {
		e.Printf("loaded profile into memory ...\n")
	}

	var (
		c *model.Confirmation = &model.Confirmation{1}
	)

	json.NewEncoder(w).Encode(&c)

	return nil
}
