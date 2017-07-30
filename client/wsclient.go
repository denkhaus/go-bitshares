package client

import (
	"net/rpc"

	"github.com/juju/errors"
	"golang.org/x/net/websocket"
)

type WSClient struct {
	client *rpc.Client
	conn   *websocket.Conn
}

func (p *WSClient) CallApi(apiID int, method string, args ...interface{}) (interface{}, error) {
	param := []interface{}{
		apiID,
		method,
		args,
	}

	var reply interface{}
	err := p.client.Call("call", param, &reply)
	return reply, err
}

//Close the client and the underlaying websocket
func (p *WSClient) Close() {
	if p.conn != nil {
		p.conn.Close()
		p.conn = nil
	}
}

//NewWebsocketClient creates a websocket based client
func NewWebsocketClient(url string) (*WSClient, error) {
	ws, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		return nil, errors.Annotate(err, "websocket dial")
	}

	cli := WSClient{
		client: NewClient(ws),
		conn:   ws,
	}

	return &cli, nil
}
