package game

type WorldNodeProperties struct {
	Name          string `json:"name"`
	Info          string `json:"info"`
	Capacity      int    `json:"capacity"`
	DeployCost    int    `json:"deployCost"`
	MoveCost      int    `json:"moveCost"`
	AttackPenalty int    `json:"attackPenalty"`
}

func generateNodeProperties(d *Data) (*WorldNodeProperties, error) {
	info, e := d.Manager.Random()
	if e != nil {
		return nil, e
	}

	p := &WorldNodeProperties{}

	p.Name = d.Manager.GenerateHyphenatedName()
	p.Info = *info
	p.Capacity = 1
	p.DeployCost = 1
	p.MoveCost = 1
	p.AttackPenalty = 1

	return p, nil
}

type WorldNodeState struct {
	IsHQ     bool `json:"isHQ"`
	Occupied bool `json:"occupied"`
	Occupant int  `json:"occupant"`
}

func initNodeState() (*WorldNodeState, error) {
	return &WorldNodeState{
		IsHQ:     false,
		Occupied: false,
		Occupant: -1,
	}, nil
}

type WorldNodeData struct {
	*WorldNodeProperties `json:"properties"`
	*WorldNodeState      `json:"state"`
}

func NewWorldNodeData() *WorldNodeData {
	return &WorldNodeData{&WorldNodeProperties{}, &WorldNodeState{}}
}

type WorldNode struct {
	LookupKey   RedisKey
	PositionKey RedisKey

	X float32
	Y float32

	*WorldNodeData
}

func NewWorldNode(lk RedisKey, pk RedisKey, x float32, y float32, d *Data) (*WorldNode, error) {
	p, e := generateNodeProperties(d)
	if e != nil {
		return nil, e
	}

	s, e := initNodeState()
	if e != nil {
		return nil, e
	}

	wn := &WorldNode{
		LookupKey:   lk,
		PositionKey: pk,

		X: x,
		Y: y,

		WorldNodeData: &WorldNodeData{
			WorldNodeProperties: p,
			WorldNodeState:      s,
		},
	}

	return wn, nil
}

func (wn *WorldNode) GetLookupKey() string   { return string(wn.LookupKey) }
func (wn *WorldNode) GetPositionKey() string { return string(wn.PositionKey) }
