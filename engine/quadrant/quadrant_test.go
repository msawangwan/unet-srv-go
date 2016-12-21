package quadrant

import (
	"math/rand"
	"testing"
)

func TestCreateNewTree(t *testing.T) {
	t.Log("create a quadrant tree with four child node slots ...")

	var (
		nodeCount  int     = 4
		nodeRadius float32 = 0.5
	)

	qt := New(nodeCount, nodeRadius)

	for i, node := range qt.Nodes {
		t.Logf("range over empty child slot %d [%s] ...", i, node)
	}

	t.Log("created a quadrant tree and child slots without error")
}

func TestPopulateQuadrantTree(t *testing.T) {
	t.Skip("skip building quad tree ...")
	t.Log("build quadrant tree ...")

	var (
		nodeCount          int     = 30
		nodeRadius         float32 = 0.5
		minScale, maxScale int     = -50, 50
		x, y               float32
	)

	qt := New(nodeCount, nodeRadius)
	r := rand.New(rand.NewSource(99))

	randInRange := func(min, max int) float32 {
		return float32(r.Intn(max-min)+min) * r.Float32()
	}

	var (
		numCreated int = 0
		attempts   int = 0
		maxAllowed int = 30
	)

	for numCreated < cap(qt.Nodes) {
		for _, v := range qt.Nodes {
			if v == nil {
				t.Log("nil node, try to insert into the quadrant tree")
				x = randInRange(minScale, maxScale)
				y = randInRange(minScale, maxScale)
				t.Logf("%f %f", x, y)
				qt.Root.tryInsert(newPoint(x, y, nodeRadius), -1)
			} else {
				if !qt.id.assigned[v.id] {
					qt.id.assigned[v.id] = true
					t.Log("created success")
					numCreated++
				}
			}
		}

		if attempts > maxAllowed {
			t.Log("max node spawn attempts reached, breaking")
			break
		}

		attempts++
	}

	t.Log("populated a quadrant tree without error")
}
