package quadrant

import (
	"testing"
)

func TestValidateIDGenerator(t *testing.T) {
	t.Log("ensure id generator doesn't get stuck in infinite loop")

	t.Log("\tinit generator...")

	var idg *idCache = &idCache{
		next:     -2,
		assigned: make(map[int]bool),
	}

	t.Log("\tgenerate ids")

	var i int = 0
	for i < 50 {
		t.Log("\t\tid: ", idg.nextID())
		i++
	}

	t.Log("\tfinished id generation")
}
