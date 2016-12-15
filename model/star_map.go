package model

import (
	"time"
)

type StarMap struct {
	Seed         int64 `json:"seed"`
	LoadExisting bool  `json:"loadExisting"`
}

func GenerateMapSeed() int64 {
	return time.Now().UTC().UnixNano()
}
