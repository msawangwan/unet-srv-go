package manager

import (
	"github.com/msawangwan/unet-srv-go/adt"
	"github.com/msawangwan/unet-srv-go/engine/prng"
)

// NameGenerator generates names that are 6 chars long, 5 digits and 1 hypen,
// where n = 0-4 numbers m = 5-n letters, an example being: 1PF-8C
type NameGenerator struct {
	*adt.StringSet
	*prng.Instance
}

// unicode constants
const (
	uHYPHEN = 45
	uZERO   = 48
	uNINE   = 57
	uA      = 65
	uZ      = 90
)

func NewNameGenerator() (*NameGenerator, error) {
	return &NameGenerator{
		StringSet: adt.NewStringSet(),
		Instance:  prng.New(0),
	}, nil
}

func (ng *NameGenerator) GenerateHyphenatedName() string {
	const length = 6

	var (
		name string
		buf  []byte
	)

	for {
		hyIndex := ng.InRangeInt(1, 5) // where is the hyphen in the name
		n := ng.InRangeInt(0, 5)       // total number of numbers in the name

		digits := 0       // numbers generated so far
		digitOrChar := -1 // 0 means random num, any other value means random letter

		buf = make([]byte, length)

		for i := 0; i < length; i++ {
			if i == hyIndex {
				buf[i] = byte(uHYPHEN)
				continue
			}

			if digits < n {
				digitOrChar = ng.InRangeInt(0, 2)
			}

			if digitOrChar == 0 {
				buf[i] = byte(ng.InRangeInt(uZERO, uNINE+1))
				digits++
			} else {
				buf[i] = byte(ng.InRangeInt(uA, uZ+1))
			}
		}

		name = string(buf)

		if ng.IsMember(name) {
			continue
		} else {
			ng.Sadd(name)
			break
		}
	}

	return name
}
