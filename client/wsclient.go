package client

import (
	"io"
	"net"

	"log"
	"sync"
	"time"

	"github.com/denkhaus/bitshares/logging"
	"github.com/juju/errors"
	"github.com/mitchellh/mapstructure"
	"github.com/pquerna/ffjson/ffjson"
	"golang.org/x/net/websocket"
)

var ErrShutdown = errors.New("connection is shut down")

type wsClient struct {
	*ffjson.Decoder
	*ffjson.Encoder
	conn        *websocket.Conn
	url         string
	resp        rpcResponse // unmarshal target
	notify      rpcNotify   // unmarshal target
	onError     ErrorFunc
	errors      chan error
	closing     bool
	shutdown    bool
	currentID   uint64
	wg          sync.WaitGroup
	mutex       sync.Mutex // protects the following
	pending     map[uint64]*RPCCall
	mutexNotify sync.Mutex // protects the following
	notifyFns   map[int]NotifyFunc
}

func NewWebsocketClient(endpointURL string) WebsocketClient {
	cli := wsClient{
		pending:   make(map[uint64]*RPCCall),
		notifyFns: make(map[int]NotifyFunc),
		currentID: 1,
		url:       endpointURL,
	}

	return &cli
}

func (p *wsClient) Connect() error {
	conn, err := websocket.Dial(p.url, "", "http://localhost/")
	if err != nil {
		return errors.Annotate(err, "Dial")
	}

	p.shutdown = false
	p.closing = false

	p.errors = make(chan error, 10)
	p.Decoder = ffjson.NewDecoder()
	p.Encoder = ffjson.NewEncoder(conn)
	p.conn = conn

	p.wg.Add(1)
	go p.monitor()

	p.wg.Add(1)
	go p.receive()

	return nil
}

func (p *wsClient) Close() error {
	if p.conn != nil {
		p.closing = true
		if err := p.conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
			return errors.Annotate(err, "SetWriteDeadline")
		}
		if err := p.conn.Close(); err != nil {
			return errors.Annotate(err, "Close [conn]")
		}

		p.wg.Wait()
		close(p.errors)
		p.conn = nil
	}

	return nil
}

func (p *wsClient) IsConnected() bool {
	if p.shutdown || p.closing {
		return false
	}

	return p.conn != nil
}

func (p *wsClient) monitor() {
	defer p.wg.Done()

	for !p.shutdown {
		select {
		case err := <-p.errors:
			if err != nil {
				if p.onError != nil {
					p.onError(err)
				} else {
					log.Println("rpc error: ", err)
				}
			}
		default:
		}
	}
}

func (p *wsClient) handleCustomData(data map[string]interface{}) error {
	logging.DumpJSON("ws notify <", data)

	switch {
	case p.notify.Is(data):
		p.notify.reset()
		err := mapstructure.Decode(data, &p.notify)
		if err != nil {
			return errors.Annotate(err, "Decode [notify]")
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
	defer p.wg.Done()

	for !p.closing {
		//TODO: is there a faster way to distinguish between RPCResponse and RPCNotify data
		var data map[string]interface{}
		if err := p.DecodeReader(p.conn, &data); err != nil {
			if e, ok := err.(*net.OpError); ok {
				if e.Err.Error() == "use of closed network connection" {
					// end loop without notification
					break
				}
			}

			//report all errors but EOF
			if err != io.EOF {
				p.errors <- errors.Annotate(err, "DecodeReader")
			}

			continue
		}

		if p.resp.Is(data) {
			p.resp.reset()
			if err := mapstructure.Decode(data, &p.resp); err != nil {
				p.errors <- errors.Annotate(err, "Decode [resp]")
				continue
			}

			logging.DumpJSON("ws resp <", data)

			if call, ok := p.pending[p.resp.ID]; ok {
				p.mutex.Lock()
				delete(p.pending, p.resp.ID)
				p.mutex.Unlock()

				call.Reply = p.resp.Result
				if p.resp.HasError() {
					call.Error = p.resp.Error
				}

				call.done()
			} else {
				p.errors <- errors.Errorf("no corresponding call found for incoming rpc data %v", p.resp)
				continue
			}
		} else if err := p.handleCustomData(data); err != nil {
			p.errors <- errors.Annotate(err, "handleCustomData")
			continue
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
		return nil, errors.Annotate(err, "Call")
	}

	<-call.Done
	return call.Reply, call.Error
}

func (p *wsClient) Call(method string, args []interface{}) (*RPCCall, error) {
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

	logging.DumpJSON("ws req >", call.Request)

	if err := p.conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return nil, errors.Annotate(err, "SetWriteDeadline")
	}

	if err := p.Encode(call.Request); err != nil {
		p.mutex.Lock()
		delete(p.pending, call.Request.ID)
		p.mutex.Unlock()

		return nil, errors.Annotate(err, "Encode [req]")
	}

	return call, nil
}
