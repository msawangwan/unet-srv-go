package quadrant

import (
	"testing"
)

func TestValidateIDGenerator(t *testing.T) {
	t.Log("ensure id generator doesn't get stuck in infinite loop")

	t.Log("\tinit id generator ...")

	var idg *idCache = &idCache{
		next:     -2,
		assigned: make(map[int]bool),
	}

	t.Log("\tgenerate ids ...")

	var i int = 0
	for i < 50 {
		t.Logf("\tid [# %d]: %d", i, idg.nextID())
		i++
	}

	t.Logf("\tfinished id generation, generated %d ids without error", i)
}
