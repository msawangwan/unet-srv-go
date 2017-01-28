package serialization

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

// TODO: currently broken and does not work
func LoadJSONFile(fname string, T interface{}) (interface{}, error) {
	if !strings.HasSuffix(fname, ".json") {
		return nil, errors.New("not a json file (info text file)")
	}

	f, e := os.Open(fname)
	if e != nil {
		return nil, e
	}

	if T == nil {
		return nil, errors.New("nil interace")
	}

	e = json.NewDecoder(f).Decode(&T)
	if e != nil {
		return nil, e
	}

	return T, nil
}
