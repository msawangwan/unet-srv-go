// package quadrant is a quadtree-like data structure, useful for creating sparse
// space/star maps for video games

package quadrant

import (
	"fmt"
)

// type Subdivider is the interface implemented by types that can sort nodes into quadrants
type Subdivider interface {
	Subdivide()
}

// type point is a 2D coordinate that covers an area defined by radius
type point struct {
	x, y, radius float32
}

func newPoint(x, y, r float32) point {
	return point{x: x, y: y, radius: r}
}

func (p point) String() string { return fmt.Sprintf("point: <%f, %f> radius: %f", p.x, p.y, p.radius) }

// type node defines the properties of a quadrant
type node struct {
	point
	id
	subquadrants []*node
	depth        int
	label        string
}

func newNode(p point, depth int, label string) *node {
	return &node{
		point:        newPoint(p.x, p.y, p.radius),
		subquadrants: make([]*node, 4),
		depth:        depth,
		label:        label,
	}
}

func (n *node) isOverlappedBy(p point) bool { return false }

func (n *node) tryInsert(p point, depth int) {
	depth++

	x1 := n.x
	y1 := n.y

	x2 := p.x
	y2 := p.y

	first := (x2 > x1) && (y2 > y1)
	second := (x2 > x1) && (y2 < y1)
	third := (x2 < x1) && (y2 < y1)
	fourth := (x2 < x1) && (y2 > y1)

	if first {
		if n.subquadrants[0] == nil {
			if !n.isOverlappedBy(p) {
				n.subquadrants[0] = newNode(p, depth, "[quadrant_1]")
			}
		} else {
			n.subquadrants[0].tryInsert(p, depth)
		}
	} else if second {
		if n.subquadrants[1] == nil {
			if !n.isOverlappedBy(p) {
				n.subquadrants[1] = newNode(p, depth, "[quadrant_2]")
			}
		} else {
			n.subquadrants[1].tryInsert(p, depth)
		}
	} else if third {
		if n.subquadrants[2] == nil {
			if !n.isOverlappedBy(p) {
				n.subquadrants[2] = newNode(p, depth, "[quadrant_3]")
			}
		} else {
			n.subquadrants[2].tryInsert(p, depth)
		}
	} else if fourth {
		if n.subquadrants[3] == nil {
			if !n.isOverlappedBy(p) {
				n.subquadrants[3] = newNode(p, depth, "[quadrant_4]")
			}
		} else {
			n.subquadrants[3].tryInsert(p, depth)
		}
	} else {
		fmt.Sprintf("unable to insert:\n\t [%s]", p)
	}
}

func (n *node) String() string {
	return fmt.Sprintf("quadrant node: [%s] id: [%d] depth: [%d] label: [%s]", n.point, n.id, n.depth, n.label)
}

// type tree consists of a root node (and it's children) that is the parent of all subquadrants
type tree struct {
	Root  *node
	Nodes []*node
	*store
}

func New(nodeCount int, nodeRadius float32) *tree {
	var (
		r *node
		s *store
	)

	s = NewIDStore(-2)

	r = newNode(newPoint(0, 0, nodeRadius), -1, "root_quadrant")
	r.id = s.nextAvailable()

	return &tree{
		Root:  r,
		Nodes: make([]*node, nodeCount),
		store: s,
	}
}

func (t *tree) AddQuadrant(n *node, i int) {
	t.Nodes[i] = n
}

func (t *tree) Subdivide(max) {
	var (
		created    []id
		a, ma      int
		smin, smax float32
	)

}

func (t *tree) String() string { return fmt.Sprintf("quadrant tree root:\n\t%v\n", t.Root) } // TODO: range over children
