package game

import "github.com/msawangwan/unet-srv-go/engine/manager"

type properties struct {
	Name          string
	info          string
	Capacity      int
	DeployCost    int
	MoveCost      int
	AttackPenalty int
}

func generateNodeProperties(ng *manager.NameGenerator) (*properties, error) {
	var (
		p *properties
	)

	p.Name = ng.GenerateHyphenatedName()

	return p, nil
}

type state struct {
	Valid    bool
	IsHQ     bool
	Occupied bool
	Occupant int
}

func initNodeState() (*state, error) {
	return &state{
		Valid:    true,
		IsHQ:     false,
		Occupied: false,
		Occupant: -1,
	}, nil
}

type WorldPositionNode struct {
	Key       RedisKey
	UniqueKey RedisKey

	X float32
	Y float32

	*properties
	*state

	*manager.NameGenerator
}

func NewWorldPositionNode(key RedisKey, uniqueKey RedisKey, x float32, y float32, ng *manager.NameGenerator) (*WorldPositionNode, error) {
	p, e := generateNodeProperties(ng)
	if e != nil {
		return nil, e
	}

	s, e := initNodeState()
	if e != nil {
		return nil, e
	}

	wpn := &WorldPositionNode{
		Key:       key,
		UniqueKey: uniqueKey,

		X: x,
		Y: y,

		properties: p,
		state:      s,
	}

	return wpn, nil
}
