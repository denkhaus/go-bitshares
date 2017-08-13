package api

import (
	"github.com/denkhaus/bitshares/client"
	"github.com/juju/errors"
)

type BitsharesApi struct {
	client *client.WSClient
}

func (p *BitsharesApi) Close() {
	if p.client != nil {
		p.client.Close()
		p.client = nil
	}
}

func New(url string) (*BitsharesApi, error) {
	client, err := client.NewWebsocketClient(url)
	if err != nil {
		return nil, errors.Annotate(err, "new websocket client")
	}

	api := BitsharesApi{
		client: client,
	}

	return &api, nil
}
