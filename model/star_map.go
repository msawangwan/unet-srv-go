package model

type StarMap struct {
	Seed         int64   `json:"seed"`
	StarCount    int     `json:"starCount"`
	Scale        float64 `json:"scale"`
	Density      float64 `json:"density"`
	LoadExisting bool    `json:"loadExisting"`
}

// see:
// https://github.com/lordofduct/spacepuppy-unity-framework/blob/master/SpacepuppyBase/Utils/RandomUtil.cs

// need:
// random.insideUnitSphere() ... returns v2 unit vector
//	func RandOnUnitSphere() {
//	a:= rand.Next() * TWO_PI
//	b := rand.Next() * TWO_PI
//	sin := Mathf.Sin(a)
//	return new Vector2(sin * mathf.Cos(b), sin * mathf.Sin(b), mathf.cos(a))
//}

//func randOnUnitCircle() {
//	a := rng.Next() * TWO_PI
//	return Vector2(mathf.sin(a), mathf.cos(b))
//}

//func insideUnitSphere() {
//	return Randonunitsphere() * rng.Next()
//}

//func insideUnitCircle() {
//	return randOnUnitCircle()
//}

// circlecast2d
