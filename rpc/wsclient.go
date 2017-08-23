package rpc

import (
	"encoding/json"

	"fmt"
	"log"
	"sync"
	"time"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/websocket"
)

var ErrShutdown = errors.New("connection is shut down")

type wsClient struct {
	*json.Decoder
	*json.Encoder
	conn        *websocket.Conn
	url         string
	resp        rpcResponse // unmarshal target
	notify      rpcNotify   // unmarshal target
	onError     ErrorFunc
	errors      chan error
	done        chan struct{}
	closing     bool
	shutdown    bool
	currentID   uint64
	mutex       sync.Mutex // protects the following
	pending     map[uint64]*RPCCall
	mutexNotify sync.Mutex // protects the following
	notifyFns   map[int]NotifyFunc
}

func NewWebsocketClient(endpointURL string) WebsocketClient {
	cli := wsClient{
		pending:   make(map[uint64]*RPCCall),
		errors:    make(chan error, 10),
		notifyFns: make(map[int]NotifyFunc),
		done:      make(chan struct{}, 1),
		currentID: 1,
		url:       endpointURL,
	}

	return &cli
}

func (p *wsClient) Connect() error {
	conn, err := websocket.Dial(p.url, "", "http://localhost/")
	if err != nil {
		return errors.Annotate(err, "dial")
	}

	p.Decoder = json.NewDecoder(conn)
	p.Encoder = json.NewEncoder(conn)
	p.conn = conn

	go p.monitor()
	go p.receive()

	return nil
}

func (p *wsClient) Close() error {
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

func (p *wsClient) monitor() {
	for {
		select {
		case err := <-p.errors:
			if err != nil {
				if p.onError != nil {
					p.onError(err)
				} else {
					log.Println("rpc error: ", err)
				}
			}

		case <-p.done:
			break
		}
	}
}

func (p *wsClient) handleCustomData(data map[string]interface{}) error {
	//	util.Dump("custom data >", data)

	switch {
	case p.notify.Is(data):
		p.notify.reset()
		err := mapstructure.Decode(data, &p.notify)
		if err != nil {
			return errors.Annotate(err, "decode notify")
		}

		params := p.notify.Params.([]interface{})
		subscriberID := int(params[0].(float64))

		var fn NotifyFunc
		p.mutexNotify.Lock()
		fn = p.notifyFns[subscriberID]
		p.mutexNotify.Unlock()

		if fn != nil {
			if err := fn(params[1]); err != nil {
				return errors.Annotate(err, "handle notify")
			}
		}
	default:
		return errors.Errorf("unhandled custom data: %v", data)
	}

	return nil
}

func (p *wsClient) receive() {

loop:
	for {
		select {
		case <-p.done:
			break loop
		default:
			//TODO: is there a faster way to distinguish between RPCResponse and RPCNotify data
			var data map[string]interface{}
			if err := p.Decode(&data); err != nil {
				util.Dump("err1", err)
				p.errors <- errors.Annotate(err, "decode in")
				break
			}

			if p.resp.Is(data) {
				p.resp.reset()
				err := mapstructure.Decode(data, &p.resp)
				if err != nil {
					util.Dump("err2", err)
					p.errors <- errors.Annotate(err, "decode response")
					break
				}

				//util.Dump(">", p.resp)

				if call, ok := p.pending[p.resp.ID]; ok {
					p.mutex.Lock()
					delete(p.pending, p.resp.ID)
					p.mutex.Unlock()

					call.Reply = p.resp.Result
					if p.resp.Error != nil {
						call.Error = formatError(p.resp.Error)
					}

					call.done()
				} else {
					p.errors <- errors.Errorf("no corresponding call found for incomming rpc data %v", p.resp)
					continue
				}
			} else if err := p.handleCustomData(data); err != nil {
				util.Dump("err3", err)
				p.errors <- errors.Annotate(err, "handle custom data")
				continue
			}
		}
	}

	// Terminate pending calls
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.shutdown = true
	for _, call := range p.pending {
		call.Error = ErrShutdown
		call.done()
	}
}

func (p *wsClient) OnNotify(subscriberID int, fn NotifyFunc) error {
	if _, ok := p.notifyFns[subscriberID]; ok {
		return errors.Errorf("a notify hook for subscriberID %d is already defined", subscriberID)
	}

	p.mutexNotify.Lock()
	p.notifyFns[subscriberID] = fn
	p.mutexNotify.Unlock()

	return nil
}

func (p *wsClient) OnError(fn ErrorFunc) {
	p.onError = fn
}

func (p *wsClient) CallAPI(apiID int, method string, args ...interface{}) (interface{}, error) {
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

func (p *wsClient) Call(method string, args interface{}) (*RPCCall, error) {
	if p.shutdown || p.closing {
		return nil, ErrShutdown
	}

	call := &RPCCall{
		Request: rpcRequest{
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

	//util.DumpJSON("rpc >", call.Request)

	if err := p.conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return nil, errors.Annotate(err, "set write deadline")
	}

	if err := p.Encode(call.Request); err != nil {
		p.mutex.Lock()
		delete(p.pending, call.Request.ID)
		p.mutex.Unlock()

		return nil, errors.Annotate(err, "encode")
	}

	return call, nil
}

func formatError(err interface{}) error {
	if e, ok := err.(map[string]interface{}); ok {
		out, _ := json.MarshalIndent(e, "", " ")
		return fmt.Errorf("server error: %s", out)
	}

	return fmt.Errorf("server error: %s", err)
}
