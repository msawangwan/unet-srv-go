package model

import (
	"github.com/msawangwan/unitywebservice/db"
	//	"log"
)

type PlayerProfile struct {
	Name        string
	UUID        uint64
	DateCreated string
}

func GetPlayersTest() (playerprofiles []PlayerProfile, err error) {
	rows, err := data.PostGreService.Query(
		"SELECT profile_name, profile_uuid, date_created FROM player_profile ORDER BY date_created DESC",
	)

	if err != nil {
		return
	}

	for rows.Next() {
		p := PlayerProfile{}

		if err = rows.Scan(&p.Name, &p.UUID, &p.DateCreated); err != nil {
			return
		}

		playerprofiles = append(playerprofiles, p)
	}

	rows.Close()

	return
}
