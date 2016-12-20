package engine

import (
	"github.com/msawangwan/unitywebservice/mathf"
	"testing"
)

func TestNewTree(t *testing.T) {
	t.Log("create a new tree")

	v := mathf.NewVector2f(5, 13)
	t.Log("v ", v)
	t.Log("now create")
	tree := CreateNewQuadTree(v)

	t.Log("tree: ", tree)

	vrr := []mathf.Vector2f{
		mathf.NewVector2f(3, 8),
		mathf.NewVector2f(6, 4),
		mathf.NewVector2f(2, 3),
		mathf.NewVector2f(12, 7),
		mathf.NewVector2f(8, 4),
		mathf.NewVector2f(1, 2),
		mathf.NewVector2f(5, 6),
		mathf.NewVector2f(14, 3),
		mathf.NewVector2f(4, 9),
		mathf.NewVector2f(15, 1),
	}

	for _, n := range vrr {
		tree.AddNewNode(n)
	}
}
