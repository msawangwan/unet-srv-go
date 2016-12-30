package model

type GameParameter struct {
	MaximumAttemptsWhenSpawningNodes int     `json:"maximumAttemptsWhenSpawningNodes"`
	NodeCount                        int     `json:"nodeCount"`
	WorldScale                       float32 `json:"worldScale"`
	NodeRadius                       float32 `json:"nodeRadius"`
}

func NewGameParameter(maxAttempts int, nodeCount int, worldScale float32, nodeRadius float32) *GameParameter {
	return &GameParameter{
		MaximumAttemptsWhenSpawningNodes: maxAttempts,
		NodeCount:                        nodeCount,
		WorldScale:                       worldScale,
		NodeRadius:                       nodeRadius,
	}
}
