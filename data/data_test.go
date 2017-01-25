package data

import (
	"io"
	"testing"
)

var path = "../schema.json"

func TestDoesPrettyPrintingWork(t *testing.T) {

	s, err := MarshallSchema(path)
	if err != nil {
		if err != io.EOF {
			t.Errorf("ERR: %s", err)
		}
	}
	t.Logf("\n%s", s.PrettyPrint())
}

func TestGetASpecifictype(t *testing.T) {
	s, err := MarshallSchema(path)
	if err != nil {
		if err != io.EOF {
			t.Errorf("ERR:: %s", err)
		}
	}

	for k, v := range s.Schemas {
		if k == "world_position_node" {
			t.Logf("\n%v", v)
		}
	}
}
