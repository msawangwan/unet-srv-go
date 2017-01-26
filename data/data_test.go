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
			t.Errorf("ERR: %s", err)
		}
	}

	schema, err := s.Access("world_position_node")
	if err != nil {
		t.Errorf("ERR: %s", err)
	}

	t.Logf("%v", schema)
}
