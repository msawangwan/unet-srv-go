package service

import (
	"encoding/json"
	"github.com/msawangwan/unitywebservice/model"
	"github.com/msawangwan/unitywebservice/util"
	//	"github.com/msawangwan/unitywebservice/db"
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
	}

	err = profile.MarkNameAsNotAvailable()
	if err != nil {
		util.Log.DbErr(w, r, err)
	} else {
		util.Log.DbActivity("appended " + profile.Name + " to list of unavailable names")
	}

	json.NewEncoder(w).Encode(profile)
}

/* expects a 'starmap' struct */
func starMap(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "nil req body", 400)
		return
	}

	var sm model.StarMap

	err := json.NewDecoder(r.Body).Decode(&sm)
	if err != nil {
		http.Error(w, "error decoding json "+err.Error(), 400)
		return
	}

	if sm.LoadExisting {
		//sm.Seed = model.GenerateMapSeed()
		log.Printf("NOT IMPLEMENTED\n")
	} else {
		sm.Seed = model.GenerateMapSeed()
	}

	json.NewEncoder(w).Encode(&sm)

	log.Printf("responded to request with a new state seed: %+v\n", sm)
}
