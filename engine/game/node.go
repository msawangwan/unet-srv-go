package game

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

func ValidateNodePosition(x, y float32) (*bool, error) {
	var (
		isValid *bool
	)

	return isValid, nil
}

type WorldNode struct{}

func (wn *WorldNode) GenerateStats(p *pool.Pool, l *debug.Log) error {
	return nil
}
