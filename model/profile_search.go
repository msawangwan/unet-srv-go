package model

import (
	"github.com/msawangwan/unet/db"
	"github.com/msawangwan/unet/util"
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
		util.Log.DbActivity("there's already a player named: " + ps.Name)
		ps.IsAvailable = false
		return false, nil
	} else {
		util.Log.DbActivity("name is available for use: " + ps.Name)
		ps.IsAvailable = true
		return true, nil
	}
}
