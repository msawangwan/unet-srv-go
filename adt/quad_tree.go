package adt

import (
	"fmt"
	"github.com/msawangwan/unitywebservice/mathf"
)

type Node struct {
	Links []*Node
	mathf.Vector2f
	depth string
}

func CreateNewQuadTree(v mathf.Vector2f) *Node {
	return &Node{
		Links:    make([]*Node, 4),
		Vector2f: mathf.NewVector2f(v.X, v.Y),
	}
}

func (n *Node) AddNewNode(v mathf.Vector2f) {
	node := &Node{
		Links:    make([]*Node, 4),
		Vector2f: mathf.NewVector2f(v.X, v.Y),
	}

	fmt.Printf("enter a new node %v\n", node)

	n.Insert(node)
}

func (n *Node) Insert(other *Node) {
	other.depth = "\tx" + other.depth
	if other.X > n.X && other.Y > n.Y {
		// quad 1
		fmt.Printf("%s quad 1\n", other.depth)
		if n.Links[0] == nil {
			n.Links[0] = other
		} else {
			n.Links[0].Insert(other)
		}
	} else if other.X > n.X && other.Y < n.Y {
		// quad 2
		fmt.Printf("%s quad 2\n", other.depth)
		if n.Links[1] == nil {
			n.Links[1] = other
		} else {
			n.Links[1].Insert(other)
		}
	} else if other.X < n.X && other.Y < n.Y {
		// quad 3
		fmt.Printf("%s quad 3\n", other.depth)
		if n.Links[2] == nil {
			n.Links[2] = other
		} else {
			n.Links[2].Insert(other)
		}
	} else if other.X < n.X && other.Y > n.Y {
		// quad 4
		fmt.Printf("%s quad 4\n", other.depth)
		if n.Links[3] == nil {
			n.Links[3] = other
		} else {
			n.Links[3].Insert(other)
		}
	}
}

func (n *Node) String() string {
	return fmt.Sprintf("%v %v", n.X, n.Y)
}
