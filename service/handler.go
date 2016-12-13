package service

import (
	"encoding/json"
	"github.com/msawangwan/unitywebservice/model"
	"log"
	"net/http"
)

func availability(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "nil req body", 400)
		return
	}

	var ps model.ProfileSearch

	err := json.NewDecoder(r.Body).Decode(&ps)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	isUnique, err := ps.IsProfileNameAvailable()
	if err != nil {
		http.Error(w, err.Error(), 400)
	} else {
		if isUnique {
			ps.IsAvailable = true
		} else {
			ps.IsAvailable = false
		}
	}

	json.NewEncoder(w).Encode(ps)

	log.Printf("responded to request with answer %v", isUnique)
}
