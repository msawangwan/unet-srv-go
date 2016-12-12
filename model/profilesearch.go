package model

type ProfileSearch struct {
	Name        string `json: "name"`
	IsAvailable bool   `json: "isAvailable"`
}
