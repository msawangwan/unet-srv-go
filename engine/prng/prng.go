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

package prgn

/*
 * m, the modulus:                m > 0
 * a, the multiplier:             0 <= a < m
 * c, the increment:              0 <= c < m
 * X (sub 0), the starting value: 0 <= X < m
 */

import (
	"time"
)

const (
	A = 8253729
	C = 2396403
	M = 32767
)

type instance struct {
	seed uint64
}

// New() returns a new prn generator, pass in 0 and it will init itself
func New(s int64) *instance {
	if s == 0 {
		s = time.Now().UTC().UnixNano()
	}
	return &instance{
		seed: uint64(s),
	}
}

func (in *instance) lcg() uint64 {
	in.seed = (A*in.seed + C) % M
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

func (in *instance) Float32() float32 {
	return float32(in.lcg()) / float32(M)
}
