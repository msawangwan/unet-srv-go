package model

import (
	"github.com/msawangwan/unitywebservice/db"
	"time"
)

type Profile struct {
	Name string `json: "name"`
	UUID string `json: "uuid"`

	Seed int64 `json: "hashedgamestate"`

	DateCreated  time.Time `json: "datecreated"`
	TimeLastSave time.Time `json: "timelastsave"`
}

var (
	STMT_SEL_ALLPROFILE    string = "SELECT * FROM profile"
	STMT_INS_CREATEPROFILE string = "INSERT INTO profile (profile_name, profile_uuid, hashed_gamestate, date_created, timeof_lastsave) VALUES ($1, $2, $3, $4, $5)"
)

func CreateNewProfile(name string) (*Profile, error) {
	t0 := time.Now()

	profile := &Profile{
		Name:         name,
		UUID:         db.CreateUUID(),
		Seed:         GenerateWorldSeedValue(),
		DateCreated:  t0,
		TimeLastSave: t0,
	}

	stmt, err := db.Postgres.DB.Prepare(STMT_INS_CREATEPROFILE)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	exec, err := stmt.Exec(
		profile.Name,
		profile.UUID,
		profile.Seed,
		profile.DateCreated,
		profile.TimeLastSave,
	)

	if err != nil || exec == nil {
		return nil, err
	}

	return profile, nil
}

func (p *Profile) MarkNameAsNotAvailable() error {
	r := db.Redis.DB.Cmd(db.CMD_SADDMEM, db.K_NAMES_NOT_AVAIL, p.Name)
	if r.Err != nil {
		return r.Err
	} else {
		return nil
	}
}
