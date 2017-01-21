package quadrant

import (
	"fmt"
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
	t.Log("partition children of a quad tree into subquadrants ...")

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

func TestConvertNodeToRedisKey(t *testing.T) {
	t.Log("convert the x and y values of a node into a redis key ...")

	qt := New(9, testNodeRadius, testSeed)
	qt.Partition(testQuadrantScale, 20)

	for _, n := range qt.Nodes {
		x, y := n.Position()
		a := fmt.Sprintf("%f:%f", x, y)
		trunc := fmt.Sprintf("%.2f:%.2f", x, y)
		t.Logf("actual [%s] formatted [%s] truncated [%s]", a, trunc, n.AsRedisKey())
	}

	t.Log("test complete")
}
