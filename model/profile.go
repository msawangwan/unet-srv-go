package model

import (
	"time"
)

type Profile struct {
	Name string `json: "name"`
	UUID string `json: "uuid"`

	HashedGameState int64 `json: "hashedgamestate"`

	DateCreated  time.Time `json: "datecreated"`
	TimeLastSave time.Time `json: "timelastsave"`
}
