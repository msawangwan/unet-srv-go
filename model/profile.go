package model

import (
	"github.com/msawangwan/unet/db"
	"time"
)

type Profile struct {
	Name string `json: "name"`
	UUID string `json: "uuid"`

	Seed int64 `json: "hashedgamestate"`

	DateCreated  time.Time `json: "datecreated"`
	TimeLastSave time.Time `json: "timelastsave"`
}

func CreateNewProfile(name string, postgre *db.PostgreHandle) (*Profile, error) {
	t0 := time.Now()

	profile := &Profile{
		Name:         name,
		UUID:         db.CreateUUID(),
		Seed:         GenerateWorldSeedValue(),
		DateCreated:  t0,
		TimeLastSave: t0,
	}

	stmt, err := postgre.Prepare(db.STATEMENT_INSERT_CREATE_PROFILE)
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

func (p *Profile) MarkNameAsNotAvailable(redis *db.RedisHandle) error {
	r := redis.Cmd(db.CMD_SADD, db.KEY_NAMES_TAKEN, p.Name)
	if r.Err != nil {
		return r.Err
	} else {
		return nil
	}
}
