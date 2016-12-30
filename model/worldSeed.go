package model

import (
	"fmt"
	"time"
)

type WorldSeed struct {
	Value int64 `json:"value"`
	IsNew bool  `json:"isNew"`
}

func (ws *WorldSeed) GenerateNew() {
	ws.Value = GenerateWorldSeedValue()
	ws.IsNew = true
}

func (ws *WorldSeed) String() string {
	return fmt.Sprintf("WorldSeed value: %v", ws.Value)
}

func CreateNewWorldSeed() *WorldSeed {
	return &WorldSeed{
		Value: GenerateWorldSeedValue(),
		IsNew: true,
	}
}

func GenerateWorldSeedValue() int64 {
	return time.Now().UTC().UnixNano()
}
