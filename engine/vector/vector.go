package vector

type vector2 struct {
	x, y float32
}

func NewVector2(x, y float32) vector2 {
	return &vector32{
		x: x,
		y: y,
	}
}
