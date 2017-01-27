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
	return &properties{}, nil
}

type state struct {
	Valid    bool
	IsHQ     bool
	Occupied bool
	Occupant int
}

func initiNodeState() (*state, error) {
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
	wpn := &WorldPositionNode{
		Key:       key,
		UniqueKey: uniqueKey,

		X: x,
		Y: y,

		properties: generateNodeProperties(),
		state:      initNodeState(),
	}

	return wpn, nil
}
