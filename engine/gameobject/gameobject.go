package gameobject

type GameObject interface {
	Instantiate() *GameObject
}
