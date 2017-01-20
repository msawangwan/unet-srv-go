package game

import (
	"fmt"
	"strconv"

	"github.com/msawangwan/unet-srv-go/engine/quadrant"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

func GameNodeString(gamelookupstr string) string {
	return fmt.Sprintf("%s:%s", gamelookupstr, "nodes")
}

func LoadWorld(gameid int, nNodes int, scale float32, nRad float32, maxA int, p *pool.Pool, l *debug.Log) error {
	conn, err := p.Get()
	if err != nil {
		return err
	}

	defer func() {
		p.Put(conn)
		l.PrefixReset()
	}()

	gamelookupstr := GameLookupString(gameid)
	gameparamstr := GameParamString(gamelookupstr)
	gamenodestr := GameNodeString(gamelookupstr)

	seedp, err := GetSeed(gameid, p, l)
	if err != nil {
		return err
	}

	seed := *seedp

	world := quadrant.New(nNodes, nRad, seed)
	world.Partition(scale, maxA)

	l.Prefix("game", "world", "load")

	nodecountstr := strconv.Itoa(nNodes)
	spawnattemptstr := strconv.Itoa(maxA)
	noderadstr := strconv.FormatFloat(float64(nRad), 'f', -1, 32)
	worldscalestr := strconv.FormatFloat(float64(scale), 'f', -1, 32)

	if err = conn.Cmd("MULTI").Err; err != nil {
		return err
	}

	for _, n := range world.Nodes {
		if !n.IsAttachedToTree() {
			l.Printf("error, detached node [%s]", n.String())
		} else {
			x, y := n.Position()
			nodestring := fmt.Sprintf("%f:%f", x, y)
			l.Printf("adding a node: [%s] [%s]\n", n.String(), nodestring)
			conn.Cmd("SADD", gamenodestr, nodestring)
		}
	}

	conn.Cmd("HMSET", gameparamstr, "node_count", nodecountstr, "world_scale", worldscalestr, "node_radius", noderadstr, "max_spawn_attempts", spawnattemptstr)

	if err = conn.Cmd("EXEC").Err; err != nil {
		return err
	}

	return nil
}

// TODO: don't need this??? was before we were getting all params
func GetSeed(gameid int, p *pool.Pool, l *debug.Log) (*int64, error) {
	conn, err := p.Get()
	if err != nil {
		return nil, err
	}

	defer func() {
		p.Put(conn)
	}()

	gameparamstr := GameParamString(GameLookupString(gameid))

	seedstring, err := conn.Cmd("HGET", gameparamstr, "world_seed").Str()
	if err != nil {
		return nil, err
	}

	seed, err := strconv.ParseInt(seedstring, 10, 64)
	if err != nil {
		return nil, err
	}

	return &seed, nil
}

type GameParameters struct {
	NodeCount        int     `json:"nodeCount"`
	NodeSpawnAttempt int     `json:"nodeMaxSpawnAttempts"`
	NodeRadius       float32 `json:"nodeRadius"`
	WorldScale       float32 `json:"worldScale"`
	Seed             int64   `json:"worldSeed"`
}

func GetGameParameters(gameid int, p *pool.Pool, l *debug.Log) (*GameParameters, error) {
	var params *GameParameters = &GameParameters{}

	// fields: node_count, world_scale, node_radius, max_spawn_attempts, world_seed
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
