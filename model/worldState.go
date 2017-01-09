package model

import (
	"fmt"
	"github.com/msawangwan/unet-srv-go/db"
	"github.com/msawangwan/unet-srv-go/engine/quadrant"
)

type WorldState struct {
	nodes *quadrant.Tree
}

func NewWorldState(sm *StarMap) *WorldState {
	return &WorldState{
		nodes: quadrant.New(sm.StarCount, sm.StarRadius, sm.Seed),
	}
}

func (ws *WorldState) LoadWorldDataIntoMem(p *Profile, redis *db.RedisHandle) error {
	conn, err := redis.Get()
	if err != nil {
		return err
	}
	defer redis.Put(conn)

	//	var (
	//		query int
	//		key   string = db.Redis.CreateKey_IsWorldInMemory("my_name")
	//	)

	//	query, err = conn.Cmd(db.CMD_GET, key)
	return nil
}

func (ws *WorldState) String() string {
	return fmt.Sprintf("world state: %s", ws.nodes)
}
