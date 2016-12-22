package quadrant

import (
	"testing"
)

func TestValidateIDGenerator(t *testing.T) {
	t.Log("ensure id generator doesn't get stuck in infinite loop")

	var ids *store = NewIDStore(-2)

	t.Log("\tgenerate ids ...")

	var i int
	for i = 0; i < 50; i++ {
		t.Logf("\tid [# %d]: %v", i, ids.nextAvailable())
		i++
	}

	t.Logf("\tfinished id generation, generated %d ids without error", i)
}
