package game

type properties struct {
	Name          string
	Info          string
	Capacity      int
	DeployCost    int
	MoveCost      int
	AttackPenalty int
}

func generateNodeProperties(d *Data) (*properties, error) {
	info, e := d.Manager.Random()
	if e != nil {
		return nil, e
	}

	p := &properties{}

	p.Name = d.Manager.GenerateHyphenatedName()
	p.Info = *info
	p.Capacity = 1
	p.DeployCost = 1
	p.MoveCost = 1
	p.AttackPenalty = 1

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
	LookupKey   RedisKey
	PositionKey RedisKey

	X float32
	Y float32

	*properties
	*state
}

func NewWorldPositionNode(lk RedisKey, pk RedisKey, x float32, y float32, d *Data) (*WorldPositionNode, error) {
	p, e := generateNodeProperties(d)
	if e != nil {
		return nil, e
	}

	s, e := initNodeState()
	if e != nil {
		return nil, e
	}

	wpn := &WorldPositionNode{
		LookupKey:   lk,
		PositionKey: pk,

		X: x,
		Y: y,

		properties: p,
		state:      s,
	}

	return wpn, nil
}

func (wpn *WorldPositionNode) GetLookupKey() string   { return string(wpn.LookupKey) }
func (wpn *WorldPositionNode) GetPositionKey() string { return string(wpn.PositionKey) }
