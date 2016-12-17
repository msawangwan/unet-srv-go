package mathf

import (
	"testing"
)

func TestRandomPointOnUnitCircle(t *testing.T) {
	var vRand Vector2f
	var s int64 = 12345

	r := NewSeededState(s)

	i := 0
	for i < 10 {
		vRand = RandOnUnitCircle(r)
		t.Log("random on unit circle: ", vRand)
		i += 1
	}
}

func TestRandomPointInsideUnitCirce(t *testing.T) {
	var vRand Vector2f
	var s int64 = 12345

	r := NewSeededState(s)

	i := 0
	for i < 10 {
		vRand = RandInsideUnitCircle(r)
		t.Log("random inside unit circle: ", vRand)
		i += 1
	}
}
