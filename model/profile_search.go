package model

import (
	"github.com/msawangwan/unet/db"
)

type ProfileSearch struct {
	Name        string `json: "name"`
	IsAvailable bool   `json: "isAvailable"`
}

func (ps *ProfileSearch) IsProfileNameAvailable(redis *db.RedisHandle) (bool, error) {
	conn, err := redis.Get()
	if err != nil {
		return false, err
	}
	defer redis.Put(conn)

	var query int

	query, err = conn.Cmd(db.CMD_SISMEMBER, db.KEY_NAMES_TAKEN, ps.Name).Int()
	if err != nil {
		return false, err
	} else {
		if query == 1 {
			ps.IsAvailable = false
			return false, nil
		} else {
			ps.IsAvailable = true
			return true, nil
		}
	}
}
