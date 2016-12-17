package mathf

import (
	"math"
	"math/rand"
)

const (
	kTWO_PI = math.Pi * 2
)

func NewSeededState(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

func RandOnUnitCircle(r *rand.Rand) Vector2f {
	var xy, x, y float64

	xy = r.Float64() * kTWO_PI

	x = math.Sin(xy)
	y = math.Cos(xy)

	return Vector2f{
		X: x,
		Y: y,
	}
}

func RandInsideUnitCircle(r *rand.Rand) Vector2f {
	return RandOnUnitCircle(r).Scale(r.Float64())
}
