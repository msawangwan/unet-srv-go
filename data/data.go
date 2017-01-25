package data

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"encoding/json"
)

type Type string
type Name string
type Properties []string

const (
	type_unit   Type = "unit"
	type_player Type = "player"
)

type NodeData struct {
	Datatype       Type       `json:"type"`
	Dataname       Name       `json:"name"`
	Dataproperties Properties `json:"properties"`
}

func LoadGameDataFile(fn string) error {
	if !strings.HasSuffix(fn, ".json") {
		return errors.New("data file not a json file")
	}

	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	var nd NodeData

	err = json.NewDecoder(f).Decode(&nd)
	if err != nil {
		return err
	}

	return nil
}

func debugPrintDataFileByLine(fn string) error {
	if !strings.HasSuffix(fn, ".json") {
		return errors.New("data file not a json file")
	}

	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		fmt.Println(s.Text())
	}

	if err = s.Err(); err != nil {
		if err != io.EOF {
			return err
		}
	}

	return nil
}
