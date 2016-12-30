package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type GameParameters struct {
	MaximumAttemptsWhenSpawningNodes int     `json:"maximumAttemptsWhenSpawningNodes"`
	WorldNodeCount                   int     `json:"worldNodeCount"`
	WorldScale                       float32 `json:"worldScale"`
	NodeRadius                       float32 `json:"nodeRadius"`
}

func LoadGameParameters(filename string) (*GameParameters, error) {
	if !strings.HasSuffix(filename, ".json") {
		return nil, errors.New("config: game parameters filename doesn't use a .json suffix")
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var (
		params *GameParameters
	)

	err = json.NewDecoder(file).Decode(&params)
	if err != nil {
		return nil, err
	}

	return params, nil
}
