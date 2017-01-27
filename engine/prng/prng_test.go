package prng

import (
	"testing"
)

const test_seed = 1482284596187742126

type iteration struct {
	i, count, times int
}

func TestPrintSomeConstants(t *testing.T) {
	t.Log("signed32: ", kMAX_SIGNED32)
	t.Log("unsigned32: ", kMAX_UNSIGNED32)
	t.Log("signed64: ", kMAX_SIGNED64)
	t.Log("unsigned64: ", kMAX_UNSIGNED64)
}

func TestJustSeeIfWeCanGenerateSomeRandomNumbers(t *testing.T) {
	t.Log("create a prng instance...")

	r := New(test_seed)

	t.Log("generate a random uint64: ", r.lcg())
	t.Log("generate a random int: ", r.Intn(5))
	t.Log("generate a random float (0,1]", r.Float32())

	t.Log("test complete")
}

func TestCanWeSeedTheGeneratorAndGetTheSameResultsBackEachTime(t *testing.T) {
	t.Log("seed the generator with a test seed ...")
	t.Log("the test seed is: ", test_seed)

	r := New(test_seed)

	var (
		i, iter int = 0, 10
		curr    uint64
		first   []uint64 = make([]uint64, iter)
	)

	t.Log("generating first set ...")

	for i = 0; i < iter; i++ {
		curr = r.lcg()
		t.Log("generated: ", curr)
		first[i] = curr
	}

	t.Log("generated first set, now re-seeding ...")

	r = New(test_seed)

	for i = 0; i < iter; i++ {
		t.Logf("compare iteration %d...", i)
		if first[i] != r.lcg() {
			t.Error("mismatched numbers generated on second pass with same seed")
		} else {
			t.Logf("generated matching numbers")
		}
	}

	t.Log("test complete")
}

func TestFloat32AlwaysReturnsABoundedFloat32(t *testing.T) {
	t.Log("create a new generator and seed it ...")

	r := New(test_seed)
	i := &iteration{i: 0, times: 10}

	var curr float32

	for i.i = 0; i.i < i.times; i.i++ {
		curr = r.Float32()
		if curr < 0 || curr >= 1 {
			t.Errorf("generated number out of bounds: %f", curr)
		} else {
			t.Logf("generated number between 0 and 1: %f", curr)
		}
	}

	t.Log("test complete")
}

func TestIntnAlwaysReturnsABoundedInt(t *testing.T) {
	t.Log("create a new generator and seed it ...")

	r := New(test_seed)
	i := &iteration{i: 0, times: 100}

	var (
		curr, max int = 0, 100
	)

	for i.i = 0; i.i < i.times; i.i++ {
		curr = r.Intn(max)
		if curr < 0 || curr > max {
			t.Errorf("generated number out of bounds: %d", curr)
		} else {
			t.Logf("generated number between 0 and %d: %d", max, curr)
		}
	}

	t.Log("test complete")
}

func TestInRangeAlwaysReturnsBoundedFloat32BetweenMinAndMax(t *testing.T) {
	t.Log("create a new generator and seed it ...")

	r := New(test_seed)
	i := &iteration{i: 0, times: 100}

	var (
		curr, min, max float32 = 0, -100.0, 100.0
	)

	for i.i = 0; i.i < i.times; i.i++ {
		curr = r.InRange(min, max)
		if curr < min || curr > max {
			t.Errorf("generated number out of bounds: %f", curr)
		} else {
			t.Logf("generated number between %f and %f: %f", min, max, curr)
		}
	}

	t.Log("test complete")
}

func TestZeroValueInputParamter(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Logf("0 as an input parameter caused no runtime panics")
		} else {
			t.Errorf("caught a panic when 0 was used as an input parameter")
		}
	}()

	t.Log("test input of 0 as a parameter ...")

	r := New(test_seed)
	i := &iteration{i: 0, times: 10}

	for i.i = 0; i.i < i.times; i.i++ {
		t.Logf("passing 0 into Intn() ... %d", r.Intn(0))
	}

	for i.i = 0; i.i < i.times; i.i++ {
		t.Logf("passing 0 into InRange() ...  %f", r.InRange(0, 0))
	}

	t.Log("test complete")
}

func TestTrigOnUnitCircle(t *testing.T) {
	t.Log("generate random 2d coordinates")

	var (
		x, y float32
	)

	r := New(test_seed)
	i := &iteration{i: 0, times: 10}

	for i.i = 0; i.i < i.times; i.i++ {
		x, y = r.onUnitCircle()
		t.Logf("random point on unit circle: <%f, %f>", x, y)
	}

	for i.i = 0; i.i < i.times; i.i++ {
		x, y = r.InUnitCircle()
		t.Logf("random point inside unit circle: <%f, %f>", x, y)
	}

	t.Log("test complete")
}

func TestHitCounter(t *testing.T) {
	r := New(test_seed)
	i := &iteration{i: 0, times: 100}

	for i.i = 0; i.i < i.times; i.i++ {
		r.lcg()
	}

	count := r.GetNumberOfValuesGenerated()
	times := i.times

	if count == times {
		t.Logf("generated [%d] which matches [%d]", count, times)
	} else {
		t.Errorf("mismatch")
	}

	t.Log("test complete")
}

func TestIntInRange(t *testing.T) {
	t.Log("test int in range")

	var min, max int = -10, 10

	r := New(test_seed)
	i := &iteration{i: 0, times: 100}

	for i.i = 0; i.i < i.times; i.i++ {
		randint := r.InRangeInt(min, max)
		if randint >= min && randint < max {
			t.Logf("int in range [%d]", randint)
		} else {
			t.Errorf("int out of range [%d]", randint)
		}
	}

	t.Logf("test completed")
}
