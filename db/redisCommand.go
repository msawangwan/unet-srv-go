package db

// redis commands
const (
	// general
	CMD_EXISTS = "EXISTS"

	// strings
	CMD_SET = "SET"
	CMD_GET = "GET"

	// lists
	CMD_RPUSH  = "RPUSH"
	CMD_LRANGE = "LRANGE"
	CMD_LINDEX = "LINDEX"
	CMD_LPOP   = "LPOP"

	// hash
	CMD_HSET    = "HSET"
	CMD_HGET    = "HGET"
	CMD_HGETALL = "HGETALL"
	CMD_HDEL    = "HDEL"

	// sets
	CMD_SADD      = "SADD"
	CMD_SREM      = "SREM"
	CMD_SISMEMBER = "SISMEMBER"
	CMD_SMEMBERS  = "SMEMBERS"
)
