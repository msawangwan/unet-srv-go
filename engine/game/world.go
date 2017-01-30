package game

import (
	"fmt"
	"strconv"

	"github.com/msawangwan/unet-srv-go/engine/manager"
	"github.com/msawangwan/unet-srv-go/engine/prng"
	"github.com/msawangwan/unet-srv-go/engine/quadrant"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

func GameValidNodeString(gamelookupstr string) string {
	return fmt.Sprintf("%s:%s", gamelookupstr, "nodes")
}

func GameNodeStatString(gamelookupstr, nodestr string) string {
	return fmt.Sprintf("%s:%s", gamelookupstr, nodestr)
}

func LoadWorld(ch *manager.ContentHandler, gameid int, nNodes int, scale float32, nRad float32, maxA int, p *pool.Pool, l *debug.Log) error {
	conn, err := p.Get()
	if err != nil {
		return err
	}

	defer func() {
		p.Put(conn)
		l.ClearLabel()
	}()

	gamelookupstr := GameLookupString(gameid)
	gameparamstr := GameParamString(gamelookupstr)
	gamevalidnodememberstr := GameValidNodeString(gamelookupstr)

	seedp, err := GetWorldSeed(gamelookupstr, p, l)
	if err != nil {
		return err
	}

	seed := *seedp

	world := quadrant.New(nNodes, nRad, prng.New(seed)) // TODO: consider...
	world.Partition(scale, maxA)

	l.Prefix("game", "world", "load")

	nodecountstr := strconv.Itoa(nNodes)
	spawnattemptstr := strconv.Itoa(maxA)
	noderadstr := strconv.FormatFloat(float64(nRad), 'f', -1, 32)
	worldscalestr := strconv.FormatFloat(float64(scale), 'f', -1, 32)

	gamedata, err := NewData(ch)
	if err != nil {
		return err
	}

	if err = conn.Cmd("MULTI").Err; err != nil {
		return err
	}

	for _, n := range world.Nodes {
		if !n.IsAttachedToTree() {
			l.Printf("error, detached node [%s]", n.String())
		} else {
			x, y := n.Position()

			noderedisstr := n.AsRedisKey()
			nodelookupstr := GameNodeStatString(gamelookupstr, noderedisstr)

			nn, e := NewWorldNode(RedisKey(nodelookupstr), RedisKey(noderedisstr), x, y, gamedata)
			if e != nil {
				l.Printf("an error occured when creating a new world position node: %s", e)
			}

			conn.Cmd("SADD", gamevalidnodememberstr, noderedisstr)
			conn.Cmd(
				"HMSET", nn.GetLookupKey(),
				"node_key", nn.GetPositionKey(),
				"node_x", nn.X,
				"node_y", nn.Y,
				"node_ishq", "false",
				"node_isoccupied", "false",
				"node_occupant", -1,
				"node_name", nn.Name,
				"node_info", nn.Info,
				"node_cap", nn.Capacity,
				"node_deploy_cost", nn.DeployCost,
				"node_move_cost", nn.MoveCost,
				"node_atk_penalty", nn.AttackPenalty,
			)

			l.Printf("added a node: [%s]\n", noderedisstr)
		}
	}

	conn.Cmd("HMSET", gameparamstr, "node_count", nodecountstr, "world_scale", worldscalestr, "node_radius", noderadstr, "max_spawn_attempts", spawnattemptstr)

	if err = conn.Cmd("EXEC").Err; err != nil {
		return err
	}

	return nil
}

// GetWorldSeed hits the redis db and returns the world seed given a game string
func GetWorldSeed(gamelookupstr string, p *pool.Pool, l *debug.Log) (*int64, error) {
	l.Label(1, "game", "world")
	l.Printf("fetching string for [lookup string (key): %s]", gamelookupstr)
	l.ClearLabel()
	ss, e := p.Cmd("HGET", GameParamString(gamelookupstr), "world_seed").Str()
	if e != nil {
		return nil, e
	}
	ws, e := strconv.ParseInt(ss, 10, 64)
	if e != nil {
		return nil, e
	}
	return &ws, nil
}

// WorldParameters wraps values that map to redis fields: node_count,
// max_spawn_attempts, world_scale, node_radius, world_seed
type WorldParameters struct {
	NodeCount        int     `json:"nodeCount"`
	NodeSpawnAttempt int     `json:"nodeMaxSpawnAttempts"`
	NodeRadius       float32 `json:"nodeRadius"`
	WorldScale       float32 `json:"worldScale"`
	Seed             int64   `json:"worldSeed"`
}

// GetWorldParameters fetches game params and retuns a struct for json encoding
func GetWorldParameters(gameid int, p *pool.Pool, l *debug.Log) (*WorldParameters, error) {
	var params *WorldParameters = &WorldParameters{}

	m, err := p.Cmd("HGETALL", GameParamString(GameLookupString(gameid))).Map()
	if err != nil {
		return nil, err
	}

	params.NodeCount, err = strconv.Atoi(m["node_count"])
	if err != nil {
		return nil, err
	}

	params.NodeSpawnAttempt, err = strconv.Atoi(m["max_spawn_attempts"])
	if err != nil {
		return nil, err
	}

	noderad, err := strconv.ParseFloat(m["node_radius"], 32)
	if err != nil {
		return nil, err
	} else {
		params.NodeRadius = float32(noderad)
	}

	worldscale, err := strconv.ParseFloat(m["world_scale"], 32)
	if err != nil {
		return nil, err
	} else {
		params.WorldScale = float32(worldscale)
	}

	params.Seed, err = strconv.ParseInt(m["world_seed"], 10, 64)
	if err != nil {
		return nil, err
	}

	return params, nil
}

func GetNodeData(gamelookupstr string, gamenodestr string, p *pool.Pool, l *debug.Log) (*WorldNodeData, error) {
	var data *WorldNodeData = NewWorldNodeData()

	m, e := p.Cmd("HGETALL", GameNodeStatString(gamelookupstr, gamenodestr)).Map()
	if e != nil {
		return nil, e
	}

	data.IsHQ, e = strconv.ParseBool(m["node_ishq"])
	if e != nil {
		return nil, e
	}

	data.Occupied, e = strconv.ParseBool(m["node_isoccupied"])
	if e != nil {
		return nil, e
	}

	data.Occupant, e = strconv.Atoi(m["node_occupant"])
	if e != nil {
		return nil, e
	}

	data.Name = m["node_name"]
	data.Info = m["node_info"]

	data.Capacity, e = strconv.Atoi(m["node_cap"])
	if e != nil {
		return nil, e
	}

	data.DeployCost, e = strconv.Atoi(m["node_deploy_cost"])
	if e != nil {
		return nil, e
	}

	data.MoveCost, e = strconv.Atoi(m["node_move_cost"])
	if e != nil {
		return nil, e
	}

	data.AttackPenalty, e = strconv.Atoi(m["node_atk_penalty"])
	if e != nil {
		return nil, e
	}

	return data, nil
}
