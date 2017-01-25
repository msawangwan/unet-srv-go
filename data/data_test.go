package data

import (
	"io"
	"testing"
)

func TestCustomStringTypes(t *testing.T) {
	var path = "../data.json"

	err := LoadGameDataFile(path)
	if err != nil {
		//	t.Errorf("%s", err)
		if err != io.EOF {
			t.Errorf("something other than eof err: %s", err)
		}
	}
}
