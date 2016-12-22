package service

import (
	"encoding/json"
	"github.com/msawangwan/unet/model"
	"github.com/msawangwan/unet/util"
	"log"
	"net/http"
)

/* expects a 'starmap' struct */
func availability(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "nil req body", 400)
		return
	}

	var ps model.ProfileSearch

	err := json.NewDecoder(r.Body).Decode(&ps)
	if err != nil {
		http.Error(w, "error decoding json "+err.Error(), 400)
		return
	}

	_, err = ps.IsProfileNameAvailable()
	if err != nil {
		http.Error(w, "error running the search "+err.Error(), 400)
	}

	json.NewEncoder(w).Encode(&ps)
}

/* expects a 'name' struct */
func profileCreate(w http.ResponseWriter, r *http.Request) {
	if err := nilBodyErr(w, r); err != nil {
		log.Printf("%v\n", err) // should log elsewyhere
		return
	}

	var gamestate *model.GameState
	var starmap *model.StarMap
	var profile *model.Profile
	var n model.ProfileName

	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		jsonDecodeErr(w, err) // should log elsewhere
		return
	}

	profile, err = model.CreateNewProfile(n.Text)
	if err != nil {
		util.Log.DbErr(w, r, err)
		return
	} else {
		util.Log.DbActivity("created a new profile with UUID " + profile.UUID + " and name of " + profile.Name)
	}

	err = profile.MarkNameAsNotAvailable()
	if err != nil {
		util.Log.DbErr(w, r, err)
	} else {
		util.Log.DbActivity("appended " + profile.Name + " to list of unavailable names")
	}

	starmap = model.NewMapDefaultParams(profile.Seed)
	gamestate = model.NewGameState(profile, starmap)

	log.Printf("new starmap created: %+v\n", starmap)

	util.Log.DbActivity("new map created")

	json.NewEncoder(w).Encode(gamestate)
}

/* expects a 'stardata' struct */
func profileNewWorldData(w http.ResponseWriter, r *http.Request) {
	if err := nilBodyErr(w, r); err != nil {
		log.Printf("%v\n", err)
		return
	}

	var sd model.StarData

	err := json.NewDecoder(r.Body).Decode(&sd)
	if err != nil {
		jsonDecodeErr(w, err)
		return
	}

	for _, v := range sd.Points {
		log.Printf("%v\n", v)
	}

	json.NewEncoder(w).Encode(&sd)
}
