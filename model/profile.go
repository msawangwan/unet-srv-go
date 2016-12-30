package model

import (
	"fmt"
	"time"

	"github.com/msawangwan/unet/db"
	"github.com/msawangwan/unet/engine/quadrant"
	"github.com/msawangwan/unet/env"
)

type Profile struct {
	Name string `json: "name"`
	UUID string `json: "uuid"`

	Seed int64 `json: "hashedgamestate"`

	DateCreated  time.Time `json: "datecreated"`
	TimeLastSave time.Time `json: "timelastsave"`
}

func CreateNewProfile(name string, postgre *db.PostgreHandle) (*Profile, error) {
	t0 := time.Now()

	profile := &Profile{
		Name:         name,
		UUID:         db.CreateUUID(),
		Seed:         GenerateWorldSeedValue(),
		DateCreated:  t0,
		TimeLastSave: t0,
	}

	stmt, err := postgre.Prepare(db.STATEMENT_INSERT_CREATE_PROFILE)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	exec, err := stmt.Exec(
		profile.Name,
		profile.UUID,
		profile.Seed,
		profile.DateCreated,
		profile.TimeLastSave,
	)

	if err != nil || exec == nil {
		return nil, err
	}

	return profile, nil
}

func (p *Profile) MarkNameAsNotAvailable(redis *db.RedisHandle) error {
	r := redis.Cmd(db.CMD_SADD, db.KEY_NAMES_TAKEN, p.Name)
	if r.Err != nil {
		return r.Err
	} else {
		return nil
	}
}

func (p *Profile) LoadIntoMemory(e *env.Global) error {
	k1 := e.CreateKey_IsWorldInMemory(p.UUID)

	conn, err := e.Get()
	if err != nil {
		return err
	}

	defer func() {
		e.Put(conn)
		e.SetPrefix_Debug()
	}()

	e.SetPrefix_DBActivity()
	e.Printf("redis: created key [%s]\n", k1)

	var (
		res int
	)

	res, err = conn.Cmd(db.CMD_EXISTS, k1).Int()
	if err != nil {
		return err
	} else {
		if res != 1 { // not in memory
			e.Printf("redis: started transaction ...\n")
			if err = conn.Cmd(db.CMD_MULTI).Err; err != nil { // start a tx
				return err
			}

			k2 := e.CreateKey_ValidWorldNodes(p.UUID)

			e.Printf("redis: created key [%s]\n", k2)

			world := quadrant.New(30, 1.2, p.Seed) // instantiate nodes
			world.Partition(50.0)

			for i, n := range world.Nodes { // store them in redis store
				x, y := n.Position()
				err = conn.Cmd(db.CMD_HSET, k2, i, fmt.Sprintf("%f:%f", x, y)).Err
				if err != nil {
					return err
				}
			}

			err = conn.Cmd(db.CMD_SET, k1, 1).Err // mark it as loaded in mem
			if err != nil {
				return err
			}

			if err = conn.Cmd(db.CMD_EXEC).Err; err != nil { // execute the tx
				return err
			}

			e.Printf("redis: executed transaction, node data now cached in memory ...\n")
		} else {
			e.Printf("redis: node data already cached ...\n")
		}
	}

	return nil
}
