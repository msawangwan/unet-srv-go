package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type Configuration struct {
	ListenAddress      string `json:"listenAddress"`
	ServicePrefix      string `json:"servicePrefix"`
	LogFile            string `json:"logFile"`
	LogDir             string `json:"logDir"`
	GameParametersFile string `json:"gameParametersFile"`
}

func LoadConfigurationFile(filename string) (*Configuration, error) {
	if !strings.HasSuffix(filename, ".json") {
		return nil, errors.New("config: configuration filename doesn't use a .json suffix")
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var (
		config *Configuration
	)

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
