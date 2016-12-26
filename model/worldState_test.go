package model

import (
	"testing"
)

const (
	testSeed = 1482284596187742126
)

var (
	testStarMap *StarMap
)

func TestWorldState(t *testing.T) {
	t.Log("world state test")
	testStarMap := NewMapDefaultParams(testSeed)
	ws := NewWorldState(testStarMap)
	t.Log("world state:", ws)
}
