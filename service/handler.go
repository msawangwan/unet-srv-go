package service

import (
	"encoding/json"
	"github.com/msawangwan/unitywebservice/db"
	"github.com/msawangwan/unitywebservice/model"
	"log"
	"net/http"
	"time"
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

	isUnique, err := ps.IsProfileNameAvailable()
	if err != nil {
		http.Error(w, "error running the search "+err.Error(), 400)
	} else {
		if isUnique {
			ps.IsAvailable = true
		} else {
			ps.IsAvailable = false
		}
	}

	json.NewEncoder(w).Encode(&ps)

	log.Printf("responded to request with answer %v", isUnique)
}

/* expects a 'name' struct */
func profileCreate(w http.ResponseWriter, r *http.Request) {
	if err := nilBodyErr(w, r); err != nil {
		log.Printf("%v\n", err)
		return
	}

	var profile *model.Profile
	var n model.ProfileName

	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		jsonDecodeErr(w, err)
		return
	}

	name := n.Text
	uuid := db.CreateUUID()
	t := time.Now()

	profile = &model.Profile{
		Name:        name,
		UUID:        uuid,
		DateCreated: t,
	}

	log.Printf("registering new player...\n")
	result, err := db.Postgres.DB.Exec("INSERT INTO profile VALUES($1, $2, $3)", profile.Name, profile.UUID, profile.DateCreated)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	log.Printf("success...%v\n", result)

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
