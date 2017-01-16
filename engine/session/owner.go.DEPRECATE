package session

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/msawangwan/unet-srv-go/debug"
)

/*
*   this class may become useful later
*
**/

// Owner may not need this type or function as it was implemented as a hack workaround
type Owner struct {
	AttachedSession string `json:"attachedSession"`
	PlayerName      string `json:"playerName"`
	PlayerIndex     int    `json:"playerIndex"`
}

// EstablishConnection was created as a result of a need for a workaround
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
