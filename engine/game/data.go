package game

import "github.com/msawangwan/unet-srv-go/engine/manager"

type Data struct {
	Manager *manager.ContentHandler
}

func NewData(ch *manager.ContentHandler) (*Data, error) {
	return &Data{
		Manager: ch,
	}, nil
}
