package mathf

import (
	"testing"
)

func TestInstantiateVectorInstances(t *testing.T) {
	t.Log("create two vector2f instances")

	var v1 Vector2f
	var v2 Vector2f

	v1 = Vector2f{
		X: 2.3,
		Y: 5.5,
	}

	v2 = Vector2f{
		X: 1.1,
		Y: 3.2,
	}

	t.Log("created: ", v1.String())
	t.Log("created: ", v2.String())
}

func TestBasicMathOperations(t *testing.T) {
	t.Log("created two vector2f structs and test basic math operations resulting in a third vector")

	var v1 Vector2f
	var v2 Vector2f
	var v3 Vector2f

	v1 = Vector2f{
		X: 2.3,
		Y: 5.5,
	}

	v2 = Vector2f{
		X: 1.1,
		Y: 3.2,
	}

	t.Log("created: ", v1.String())
	t.Log("created: ", v2.String())

	v3 = v1.Add(v2)
	t.Log("addition results: ", v3)
	v3 = v1.Subtract(v2)
	t.Log("subtraction results: ", v3)
	v3 = v1.Multiply(v2)
	t.Log("multiplication results: ", v3)
	v3 = v1.Divide(v2)
	t.Log("division results: ", v3)
	v3 = v1.Scale(5)
	t.Log("scalar multiplication results: ", v3)
	v3 = v1.Shrink(2)
	t.Log("scalar division results: ", v3)
}

func TestDivideByZero(t *testing.T) {
	t.Log("create two vectors, the divisor has x and y components of 0")

	var v1 Vector2f
	var v2 Vector2f
	var v3 Vector2f
	var v4 Vector2f

	v1 = Vector2f{
		X: 2.5,
		Y: 5.4,
	}

	v2 = Vector2f{
		X: 0.0,
		Y: 2.1,
	}

	v3 = Vector2f{
		X: 5.0,
		Y: 0.0,
	}

	v4 = v1.Divide(v2)
	t.Log("divisor with 0 x component: ", v4)
	v4 = v1.Divide(v3)
	t.Log("divisor with 0 y component: ", v4)
	v4 = v1.Shrink(0)
	t.Log("scalar division by 0: ", v4)
}

func TestMagnitudeAndSqrMagnitude(t *testing.T) {
	t.Log("common vector operations")

	var v1 Vector2f
	var v2 Vector2f
	var res float64

	v1 = Vector2f{
		X: 6.3,
		Y: 2.3,
	}

	t.Log("vector is: ", v1)

	res = v1.Magnitude()
	t.Log("mag: ", res)
	res = v1.MagnitudeSqr()
	t.Log("sqrMag: ", res)
	v2 = v1.Normalize()
	t.Log("norm: ", v2)
	res = v1.Dotf(v2)
	t.Log("dot: ", res)
	res = v1.Crossf(v2)
	t.Log("cross: ", res)
}
