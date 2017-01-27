package game

//	"github.com/msawangwan/unet-srv-go/adt"
// TODO: keep a list of names

type properties struct {
	Name          string
	info          string
	Capacity      int
	DeployCost    int
	MoveCost      int
	AttackPenalty int
}

func generateNodeProperties() (*properties, error) {
	var (
		p *properties
	)

	return &p, nil
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
}

func NewWorldPositionNode(key RedisKey, uniqueKey RedisKey, x float32, y float32) (*WorldPositionNode, error) {
	p, e := generateNodeProperties()
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
