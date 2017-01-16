package quadrant

import (
	"testing"
)

const (
	testSeed          = 1482284596187742126
	testNodeRadius    = 1.2
	testQuadrantScale = 50
)

func TestCreateNewTree(t *testing.T) {
	t.Log("create a quadrant tree with four child node slots ...")

	var (
		nodeCount  int     = 4
		nodeRadius float32 = 0.5
	)

	qt := New(nodeCount, nodeRadius, testSeed)

	for i, node := range qt.Nodes {
		t.Logf("range over empty child slot %d [%s] ...", i, node)
	}

	t.Log("created a quadrant tree and child slots without error")
}

func TestSubdividerImplementation(t *testing.T) {
	t.Log("partition children of a quad tree into subquadrants ..")

	var (
		numNodes   int     = 9
		nodeRadius float32 = testNodeRadius
		scale      float32 = testQuadrantScale
	)

	t.Log("test seed:", testSeed)

	qt := New(numNodes, nodeRadius, testSeed)

	qt.Partition(scale, 20)

	for i, n := range qt.Nodes {
		if n.attached {
			t.Logf("node is attached at point: [%s] ...", n)
		} else {
			t.Logf("found a detached node at index: %d", i)
		}
	}

	t.Log("test compelte")
}
