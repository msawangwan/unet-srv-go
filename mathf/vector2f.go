package mathf

import (
	"fmt"
	"math"
)

const (
	kTolerance = 0.0001
)

type Vector2f struct {
	X float64
	Y float64
}

func NewVector2f(x, y float64) Vector2f {
	return Vector2f{
		X: x,
		Y: y,
	}
}

func (v Vector2f) Add(other Vector2f) Vector2f {
	return Vector2f{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v Vector2f) Subtract(other Vector2f) Vector2f {
	return Vector2f{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v Vector2f) Multiply(other Vector2f) Vector2f {
	return Vector2f{
		X: v.X * other.X,
		Y: v.Y * other.Y,
	}
}

func (v Vector2f) Divide(other Vector2f) Vector2f {
	var x, y float64

	if other.X == 0 {
		x = math.NaN()
	} else {
		x = v.X / other.X
	}

	if other.Y == 0 {
		y = math.NaN()
	} else {
		y = v.Y / other.Y
	}

	return Vector2f{
		X: x,
		Y: y,
	}
}

func (v Vector2f) Scale(s float64) Vector2f {
	return Vector2f{
		X: v.X * s,
		Y: v.Y * s,
	}
}

func (v Vector2f) Shrink(s float64) Vector2f {
	var x, y float64

	if s == 0 {
		x, y = math.NaN(), math.NaN()
	} else {
		x = v.X / s
		y = v.Y / s
	}

	return Vector2f{
		X: x,
		Y: y,
	}
}

func (v Vector2f) Magnitude() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y))
}

func (v Vector2f) MagnitudeSqr() float64 {
	return ((v.X * v.X) + (v.Y * v.Y))
}

func (v Vector2f) Normalize() Vector2f {
	var x, y, m float64

	m = math.Sqrt((v.X * v.X) + (v.Y * v.Y))

	if m <= kTolerance {
		m = 1
	}

	x = v.X / m
	y = v.Y / m

	if math.Abs(x) < kTolerance {
		x = 0
	}

	if math.Abs(y) < kTolerance {
		y = 0
	}

	return Vector2f{
		X: x,
		Y: y,
	}
}

func (v Vector2f) Dotf(other Vector2f) float64 {
	return (v.X * other.X) + (v.Y * other.Y)
}

func (v Vector2f) Crossf(other Vector2f) float64 {
	return (v.X * other.Y) - (v.Y * other.X)
}

func (v Vector2f) NormalLH() Vector2f {
	return Vector2f{
		X: v.Y * -1.0,
		Y: v.X,
	}
}

func (v Vector2f) NormalRH() Vector2f {
	return Vector2f{
		X: v.Y,
		Y: v.X * -1.0,
	}
}

func (v Vector2f) Trunc(maxLen float64) Vector2f {
	if v.Magnitude() > maxLen { // TODO: optimize away the sqr calls cause we're doing em' twice here...
		return v.Normalize().Scale(maxLen)
	}
	return v
}

func (v Vector2f) String() string {
	return fmt.Sprintf("Vector2f <%v, %v>", v.X, v.Y)
}
