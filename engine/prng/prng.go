/*
 * package prng implements a (very basic lcg, for now) pseudo-random number generator that
 * that matches a c#-client side prng, based on source code from:
 *
 * http://stackoverflow.com/questions/15500621/c-c-algorithm-to-produce-same-pseudo-random-number-sequences-from-same-seed-on
 *
 * also see:
 *
 * http://www-personal.umich.edu/~mejn/percolation/rnglong.h
 * http://www-personal.umich.edu/~mejn/percolation/rnglong.c
 *
 * and read:
 *
 * https://www.gnu.org/software/gsl/manual/html_node/Unix-random-number-generators.html#Unix-random-number-generators
 */

package prng

import (
	"math"
	"time"
)

// constants for the lcg algorithm as given by knuth
const (
	// m, the modulus where m > 0
	// a, the multiplier where 0 <= a < m
	// c, the increment where 0 <= c < m
	kM = 32767
	kA = 8253729
	kC = 2396403
)

const (
	kTWO_PI = math.Pi * 2
)

const (
	kMAX_SIGNED32   = (1 << 31) - 1
	kMAX_UNSIGNED32 = (1 << 32) - 1
	kMAX_SIGNED64   = (1 << 63) - 1
	kMAX_UNSIGNED64 = uint64(1<<64 - 1)
)

type instance struct {
	seed uint64
}

// New() returns a new prn generator, pass in 0 and it will init itself
func New(s int64) *instance {
	p := kA*uint64(s) + kC

	if s <= 0 {
		s = time.Now().UTC().UnixNano()
	} else if p >= kMAX_UNSIGNED64 {
		s /= 2
	}

	return &instance{
		seed: uint64(s),
	}
}

func (in *instance) lcg() uint64 {
	in.seed = (kA*in.seed + kC) % kM
	return in.seed
}

func (in *instance) InRange(min, max float32) float32 {
	return (float32(in.Intn(int(max)-int(min))) + min) * in.Float32()
}

func (in *instance) Intn(max int) int {
	if max == 0 {
		return 0
	}
	return int(in.lcg()) % max
}

// Float32() returns a float in the range of (0, 1]
func (in *instance) Float32() float32 {
	return float32(in.lcg()) / float32(kM)
}

func (in *instance) onUnitCircle() (x, y float32) {
	h := in.Float32() * kTWO_PI

	x = float32(math.Sin(float64(h)))
	y = float32(math.Cos(float64(h)))

	return
}

func (in *instance) InUnitCircle() (x, y float32) {
	r := in.Float32()

	x, y = in.onUnitCircle()

	x *= r
	y *= r

	return
}
