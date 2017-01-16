package game

import (
	"fmt"
	"strconv"

	"github.com/msawangwan/unet-srv-go/engine/quadrant"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

func LoadWorld(gkey int, nNodes int, scale float32, nRad float32, maxA int, p *pool.Pool, l *debug.Log) error {
	conn, err := p.Get()
	if err != nil {
		return err
	}

	defer func() {
		p.Put(conn)
		l.PrefixReset()
	}()

	game, err := conn.Cmd("HGET", hk_gameHandleKey, strconv.Itoa(gkey)).Str()
	if err != nil {
		return err
	}

	seedString, err := conn.Cmd("HGET", game, hf_seed).Str()
	if err != nil {
		return err
	}

	seed, err := strconv.ParseInt(seedString, 10, 64)
	if err != nil {
		return err
	}

	worldKey := fmt.Sprintf("%s:%s", game, "nodes")

	world := quadrant.New(nNodes, nRad, seed)
	world.Partition(scale, maxA)

	if err = conn.Cmd("MULTI").Err; err != nil {
		return err
	}

	for _, n := range world.Nodes {
		if !n.IsAttachedToTree() {
			l.Printf("error, detached node [%s]", n.String())
		} else {
			x, y := n.Position()
			nodestring := fmt.Sprintf("%f:%f", x, y)
			l.Printf("adding a node: [%s] [%s]", n.String(), nodestring)
			conn.Cmd("SADD", worldKey, nodestring)
		}
	}

	if err = conn.Cmd("EXEC").Err; err != nil {
		return err
	}

	return nil
}
