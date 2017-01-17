package model

const (
	kSTAR_COUNT  = 30
	kSTAR_RADIUS = 1.2
	kWORLD_SCALE = 20.0
	kDENSITY     = 1.5
)

type StarMap struct {
	Seed         int64   `json:"seed"`
	StarCount    int     `json:"starCount"`
	StarRadius   float32 `json:"starRadius"`
	Scale        float32 `json:"scale"`
	Density      float32 `json:"density"`
	LoadExisting bool    `json:"loadExisting"`
}

func NewMapDefaultParams(seed int64) *StarMap {
	return &StarMap{
		Seed:       seed,
		StarCount:  kSTAR_COUNT,
		StarRadius: kSTAR_RADIUS,
		Scale:      kWORLD_SCALE,
		Density:    kDENSITY,
	}
}
