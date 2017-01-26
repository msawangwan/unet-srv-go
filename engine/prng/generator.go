package prng

// Generator wraps and manages a prng instances
type Generator struct {
	rand *Instance
}

func (g *Generator) Next(max int) int {
	return g.rand.Intn(max)
}

func (g *Generator) Nextf() float32 {
	return g.rand.Float32()
}

func (g *Generator) NextInRange(min, max float32) float32 {
	return g.rand.InRange(min, max)
}
