package model

import (
	"github.com/msawangwan/unitywebservice/db"
	"time"
)

type Player struct {
	Name string `json: "name"`
	UUID string `json: "uuid"`

	HashedGameState uint64 `json: "hashedgamestate"`

	DateCreated    time.Time `json: "datecreated"`
	TimeOfLastSave time.Time `json: "timeoflastsave"`
}

func SelectAllPlayers() ([]Player, error) {
	rows, err := db.Postgres.DB.Query(
		"SELECT profile_name, profile_uuid, hashed_gamestate, date_created, timeof_lastsave FROM player ORDER BY date_created DESC",
	)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var players []Player

	for rows.Next() {
		p := Player{}

		if err = rows.Scan(&p.Name, &p.UUID, &p.HashedGameState, &p.DateCreated, &p.TimeOfLastSave); err != nil {
			return nil, err
		}

		players = append(players, p)
	}

	rows.Close()

	return players, nil
}
