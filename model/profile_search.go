package model

import (
	"github.com/msawangwan/unitywebservice/db"
	"log"
)

type ProfileSearch struct {
	Name        string `json: "name"`
	IsAvailable bool   `json: "isAvailable"`
}

func (ps *ProfileSearch) IsProfileNameAvailable() (bool, error) {
	conn, err := db.Redis.DB.Get()
	if err != nil {
		return false, err
	}
	defer db.Redis.DB.Put(conn)

	var query int

	query, err = conn.Cmd(db.CMD_SISMEM, db.K_NAMES_NOT_AVAIL, ps.Name).Int()
	if err != nil {
		return false, err
	} else if query == 1 {
		log.Printf("already have a profile with name: %s\n", ps.Name)
		return false, nil
	} else {
		log.Printf("did not find a profile with name: %s\n", ps.Name)
		return true, nil
	}
}
