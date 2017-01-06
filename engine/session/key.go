package session

type Key struct {
	BareFormat  string `json:"bareFormat"`
	RedisFormat string `json:"redisFormat"`
}
