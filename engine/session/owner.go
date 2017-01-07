package session

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet/debug"
)

type Owner struct {
	AttachedSession string `json:"attachedSession"`
	PlayerName      string `json:"playerName"`
	PlayerIndex     int    `json:"playerIndex"`
}

func EstablishConnection(sid string, ip string, p *pool.Pool, l *debug.Log) (bool, *string, error) {
	conn, err := p.Get()
	if err != nil {
		return false, nil, err
	}

	defer func() {
		p.Put(conn)
		l.SetPrefixDefault()
	}()

	l.SetPrefix("[SESSION][CONNECT] ")

	//	k := e.CreateListKey_SessionConn(e.CreateHashKey_Session(sid))
	k := "session:" + sid + ":conn"

	connected, err := conn.Cmd("LRANGE", k, 0, -1).List()
	if err != nil {
		return false, nil, err
	}

	for _, c := range connected {
		if c == ip {
			break
		}
	}

	conn.Cmd("RPUSH", k, ip)

	return true, &k, nil
}
