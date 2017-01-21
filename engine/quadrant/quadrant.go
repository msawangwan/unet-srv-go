// package quadrant is a quadtree-like data structure, useful for creating sparse
// space/star maps for video games

package quadrant

import (
	"fmt"
	"github.com/msawangwan/unet-srv-go/engine/prng"
)

// type Subdivider is the interface implemented by types that can sort nodes into quadrants
type Subdivider interface {
	Partition(scale float32)
}

// type point is a 2D coordinate that covers an area defined by radius
type point struct {
	x, y, radius float32
}

func newPoint(x, y, r float32) point {
	return point{x: x, y: y, radius: r}
}

func (p point) String() string {
	return fmt.Sprintf("point: <%.4f, %.4f> radius: %.4f", p.x, p.y, p.radius)
}

// type node defines the properties of a quadrant
type node struct {
	point
	id

	subquadrants []*node

	depth    int
	label    string
	attached bool
}

func newNode(p point, depth int, label string, isattached bool) *node {
	return &node{
		point:        newPoint(p.x, p.y, p.radius),
		subquadrants: make([]*node, 4),
		depth:        depth,
		label:        label,
		attached:     isattached,
	}
}

func (n *node) isOverlappedBy(other *node) bool {
	dx := n.x - other.x
	dy := n.y - other.y

	if dx < 0 {
		dx *= -1
	}
	if dy < 0 {
		dy *= -1
	}

	if dx <= n.radius || dy <= n.radius {
		return true
	} else {
		return false
	}
}

func (n *node) tryInsert(other *node) {
	other.depth++

	x1 := n.x
	y1 := n.y

	x2 := other.x
	y2 := other.y

	first := (x2 > x1) && (y2 > y1)
	second := (x2 > x1) && (y2 < y1)
	third := (x2 < x1) && (y2 < y1)
	fourth := (x2 < x1) && (y2 > y1)

	if first {
		if n.subquadrants[0] == nil {
			if !n.isOverlappedBy(other) {
				other.label = "[quadrant_1]"
				other.attached = true
				n.subquadrants[0] = other
			}
		} else {
			n.subquadrants[0].tryInsert(other)
		}
	} else if second {
		if n.subquadrants[1] == nil {
			if !n.isOverlappedBy(other) {
				other.label = "[quadrant_2]"
				other.attached = true
				n.subquadrants[1] = other
			}
		} else {
			n.subquadrants[1].tryInsert(other)
		}
	} else if third {
		if n.subquadrants[2] == nil {
			if !n.isOverlappedBy(other) {
				other.label = "[quadrant_3]"
				other.attached = true
				n.subquadrants[2] = other
			}
		} else {
			n.subquadrants[2].tryInsert(other)
		}
	} else if fourth {
		if n.subquadrants[3] == nil {
			if !n.isOverlappedBy(other) {
				other.label = "[quadrant_4]"
				other.attached = true
				n.subquadrants[3] = other
			}
		} else {
			n.subquadrants[3].tryInsert(other)
		}
	} else {
		fmt.Sprintf("unable to insert:\n\t [%s]", other)
	}
}

func (n *node) Position() (float32, float32) { return n.x, n.y }
func (n *node) IsAttachedToTree() bool       { return n.attached }
func (n *node) AsRedisKey() string           { return forgeRedisKey(n.x, n.y) }

func (n *node) String() string {
	return fmt.Sprintf("node [%s] id [%d] depth [%d]", n.point, n.id, n.depth)
}

// ForgeRedisKey exports forgeRedisKey()
func ForgeRedisKey(x, y float32) string { return forgeRedisKey(x, y) }

func forgeRedisKey(x, y float32) string {
	trunc := func(f float32) float32 { return float32(int(f*100)) / 100 } // this doesn't round like the fmt formatter
	return fmt.Sprintf("%.2f:%.2f", trunc(x), trunc(y))
}

// type tree consists of a root node (and it's children) that is the parent of all subquadrants
type Tree struct {
	Root  *node
	Nodes []*node // TODO: deprecate and replace with a map
	//ValidNodes map[string]bool // TODO: if this comment and the one above is here, this and the idea above are not used
	size int
	*store
	*prng.Instance
}

func New(nodeCount int, nodeRadius float32, seed int64) *Tree {
	var (
		ns   []*node
		r, n *node
		s    *store
		size int
	)

	s = newIDStore(-2)
	size = nodeCount + 1

	r = newNode(newPoint(0, 0, nodeRadius), -1, "[root_quadrant]", true)
	r.id = s.nextAvailable()

	ns = make([]*node, size)
	ns[0] = r

	for i := 1; i < size; i++ {
		n = newNode(newPoint(0, 0, nodeRadius), -1, "[detached]", false)
		n.id = s.nextAvailable()
		ns[i] = n
	}

	return &Tree{
		Root:     r,
		Nodes:    ns,
		size:     size,
		store:    s,
		Instance: prng.New(seed), // TODO: pass in an instance instead
	}
}

func (t *Tree) Partition(scale float32, amax int) {
	var (
		created    map[id]bool = make(map[id]bool)
		smin, smax float32     = -scale, scale
		a, c       int         = 0, 0 // attemptsCount, createdCount
	)

	for c < t.size {
		for _, n := range t.Nodes {
			if !n.attached {
				n.x = t.Instance.InRange(smin, smax)
				n.y = t.Instance.InRange(smin, smax)
				t.Root.tryInsert(n)
			} else {
				if !created[n.id] {
					created[n.id] = true
					c++
				}
			}
		}

		if a > amax {
			fmt.Printf("engine/quadrant: building tree, max attempts reached\n")
			break
		}

		a++
	}
}

func (t *Tree) String() string { return fmt.Sprintf("quadrant tree root:\n\t%v\n", t.Root) } // TODO: range over children
