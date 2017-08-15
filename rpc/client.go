package rpc

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/websocket"
)

var ErrShutdown = errors.New("connection is shut down")

type WebsocketClient interface {
	OnError(fn func(error))
	Call(method string, args interface{}) (*RPCCall, error)
	CallApi(apiID int, method string, args ...interface{}) (interface{}, error)
	Close() error
	Connect() error
}

type RPCCall struct {
	Method  string
	Request RPCRequest
	Reply   interface{}
	Error   error
	Done    chan *RPCCall
}

func (call *RPCCall) done() {
	select {
	case call.Done <- call:
		// ok
	default:
		log.Println("rpc: discarding Call reply due to insufficient Done chan capacity")
	}
}

type RPCRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
	ID     uint64      `json:"id"`
}

type RPCResponse struct {
	ID     uint64      `json:"id"`
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

type websocketClient struct {
	conn      *websocket.Conn
	url       string
	dec       *json.Decoder
	enc       *json.Encoder
	onError   func(error)
	errors    chan error
	done      chan struct{}
	closing   bool
	shutdown  bool
	mutex     sync.Mutex // protects the following
	currentID uint64
	pending   map[uint64]*RPCCall
}

func NewWebsocketClient(url string) WebsocketClient {
	cli := websocketClient{
		pending:   make(map[uint64]*RPCCall),
		errors:    make(chan error, 10),
		done:      make(chan struct{}, 1),
		currentID: 1,
		url:       url,
	}

	return &cli
}

func (p *websocketClient) Connect() error {
	conn, err := websocket.Dial(p.url, "", "http://localhost/")
	if err != nil {
		return errors.Annotate(err, "dial")
	}

	p.dec = json.NewDecoder(conn)
	p.enc = json.NewEncoder(conn)
	p.conn = conn

	go p.monitor()
	go p.receive()

	return nil
}

func (p *websocketClient) Close() error {
	if p.conn != nil {
		p.closing = true
		if err := p.conn.Close(); err != nil {
			return errors.Annotate(err, "close connection")
		}

		p.done <- struct{}{}
		close(p.errors)
		p.conn = nil
	}

	return nil
}

func (p *websocketClient) monitor() {
	for {
		select {
		case err := <-p.errors:
			if p.onError != nil {
				p.onError(err)
			} else {
				log.Println("rpc error:  ", err)
			}
		case <-p.done:
			break
		}
	}
}

func (p *websocketClient) handleCustomData(data map[string]interface{}) error {
	util.Dump("custom >", data)
	return nil
}

func (p *websocketClient) receive() {

	for {
		var data map[string]interface{}
		if err := p.dec.Decode(&data); err != nil {
			p.errors <- errors.Annotate(err, "decode in")
			break
		}

		var resp RPCResponse
		if isRPCResponse(data) {
			err := mapstructure.Decode(data, &resp)
			if err != nil {
				p.errors <- errors.Annotate(err, "decode response")
				break
			}

			//	util.Dump(">", resp)

			if call, ok := p.pending[resp.ID]; ok {
				call.Reply = resp.Result

				if resp.Error != nil {
					call.Error = formatError(resp.Error)
				}

				call.done()
				delete(p.pending, resp.ID)
			} else {
				log.Println("unhandled response:", resp)
			}
		} else {
			if err := p.handleCustomData(data); err != nil {
				p.errors <- errors.Annotate(err, "handle custom data")
				continue
			}
		}
	}

	// Terminate pending calls.

	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.shutdown = true
	for _, call := range p.pending {
		call.Error = ErrShutdown
		call.done()
	}

}

func (p *websocketClient) OnError(fn func(error)) {
	p.onError = fn
}

func (p *websocketClient) CallApi(apiID int, method string, args ...interface{}) (interface{}, error) {
	param := []interface{}{
		apiID,
		method,
		args,
	}

	call, err := p.Call("call", param)
	if err != nil {
		return nil, errors.Annotate(err, "call")
	}

	<-call.Done
	return call.Reply, call.Error
}

func (p *websocketClient) Call(method string, args interface{}) (*RPCCall, error) {
	if p.shutdown || p.closing {
		return nil, ErrShutdown
	}

	call := &RPCCall{
		Request: RPCRequest{
			Method: method,
			Params: args,
			ID:     p.currentID,
		},
		Done: make(chan *RPCCall, 10),
	}

	p.mutex.Lock()
	p.currentID++
	p.pending[call.Request.ID] = call
	p.mutex.Unlock()

	//util.Dump(">", call.Request)

	if err := p.enc.Encode(call.Request); err != nil {
		p.mutex.Lock()
		delete(p.pending, call.Request.ID)
		p.mutex.Unlock()

		return nil, errors.Annotate(err, "encode")
	}

	return call, nil
}

func formatError(err interface{}) error {
	e, ok := err.(string)
	if !ok {
		return fmt.Errorf("invalid error %v", err)
	}

	if e == "" {
		e = "unspecified"
	}

	return fmt.Errorf("server error: %s", e)
}

func isRPCResponse(data map[string]interface{}) bool {
	if _, ok := data["id"]; ok {
		if _, ok := data["result"]; ok {
			return true
		}
	}
	return false
}
