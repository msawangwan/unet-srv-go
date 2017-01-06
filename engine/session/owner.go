package session

import (
	"github.com/msawangwan/unet/env"
)

type Owner struct {
	AttachedSession string `json:"attachedSession"`
	PlayerName      string `json:"playerName"`
	PlayerIndex     int    `json:"playerIndex"`
}

func EstablishConnection(e *env.Global, sid string, ip string) (bool, *string, error) {
	conn, err := e.Get()
	if err != nil {
		return false, nil, err
	}

	defer func() {
		e.Put(conn)
		e.SetPrefixDefault()
	}()

	e.SetPrefix("[SESSION][CONNECT] ")

	k := e.CreateListKey_SessionConn(e.CreateHashKey_Session(sid))

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
