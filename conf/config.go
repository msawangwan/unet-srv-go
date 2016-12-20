package config

import (
	"json"
	"log"
	"os"
	"strings"
)

type Configuration struct {
	ServicePrefix string
	Logfile       string
	LogDir        string
}

func NewServiceConfig(filename string) *Configuration {
	if !strings.HasSuffix(filename, ".json") {
		log.Fatalln("error: configuration isn't a json file!")
	}

	if file, err := os.Open(filename); err != nil {
		log.Fatalln("error: problems opening configuration file", err)
	}

	var cfg Configuration

	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		log.Fatalln("error: couldin't decode configuration json", err)
	}

	return cfg
}
