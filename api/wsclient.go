package api

import (
	"encoding/json"
	"io"
	"math"
	"net"
	"os"
	"syscall"

	"sync"
	"time"

	"github.com/denkhaus/logging"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/tevino/abool"
	"golang.org/x/net/websocket"
)

var (
	DialerTimeout    = time.Duration(5 * time.Second)
	ReadWriteTimeout = time.Duration(10 * time.Second)
	ErrShutdown      = errors.New("connection is shut down")
)

type wsClient struct {
	url            string
	onError        ErrorFunc
	errors         chan error
	conn           *websocket.Conn
	closing        *abool.AtomicBool
	shutdown       *abool.AtomicBool
	requestID      uint64
	subscribeID    uint64
	wg             sync.WaitGroup
	mutex          sync.Mutex // protects the following
	pending        map[uint64]*RPCCall
	mutexSubscribe sync.Mutex // protects the following
	subscrFns      map[uint64]SubscribeCallback
}

func NewWebsocketClient(endpointURL string) WebsocketClient {
	cli := wsClient{
		closing:   abool.NewBool(false),
		shutdown:  abool.NewBool(false),
		pending:   make(map[uint64]*RPCCall),
		subscrFns: make(map[uint64]SubscribeCallback),
		url:       endpointURL,
	}

	return &cli
}

func (p *wsClient) Connect() error {
	config, err := websocket.NewConfig(p.url, "http://localhost/")
	if err != nil {
		return errors.Annotate(err, "NewConfig")
	}

	config.Dialer = &net.Dialer{
		Timeout: DialerTimeout,
	}

	conn, err := websocket.DialConfig(config)
	if err != nil {
		return errors.Annotate(err, "Dial")
	}

	p.shutdown.UnSet()
	p.closing.UnSet()

	p.errors = make(chan error, 10)
	p.conn = conn

	p.wg.Add(1)
	go p.monitor()

	p.wg.Add(1)
	go p.receive()

	return nil
}

func (p *wsClient) Close() error {
	if p.conn != nil {
		p.closing.Set()
		if !p.shutdown.IsSet() {
			if err := p.conn.SetDeadline(time.Now().Add(ReadWriteTimeout)); err != nil {
				return errors.Annotate(err, "SetDeadline")
			}
			if err := p.conn.Close(); err != nil {
				return errors.Annotate(err, "Close [conn]")
			}
		}

		p.wg.Wait()
		close(p.errors)
		p.conn = nil
	}

	return nil
}

func (p *wsClient) IsConnected() bool {
	if p.shutdown.IsSet() || p.closing.IsSet() {
		return false
	}
	return p.conn != nil
}

func (p *wsClient) monitor() {
	defer p.wg.Done()

	for !p.shutdown.IsSet() {
		select {
		case err := <-p.errors:
			if err != nil {
				if p.onError != nil {
					p.onError(err)
				} else {
					logging.Errorf("WebsocketClient error: %s", err)
					logging.Warn("please set the API OnError hook to avoid this message")
				}
			}
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func (p *wsClient) mustEndReceive(err error) bool {
	if e, ok := err.(*net.OpError); ok {
		if e.Err.Error() == "use of closed network connection" {
			return true
		}

		if syscallErr, ok := e.Err.(*os.SyscallError); ok {
			if syscallErr.Err == syscall.ECONNRESET {
				return true
			}
		}
	}

	if e, ok := err.(net.Error); ok {
		if e.Timeout() {
			return true
		}
	}

	return false
}

func (p *wsClient) receive() {
	defer p.wg.Done()

	for !p.closing.IsSet() {
		var data string
		if err := websocket.Message.Receive(p.conn, &data); err != nil {
			if p.mustEndReceive(err) {
				break
			}

			//report all errors but EOF
			if err != io.EOF {
				p.errors <- errors.Annotate(err, "Receive")
			}
			continue
		}

		var resp rpcResponse
		if err := ffjson.Unmarshal([]byte(data), &resp); err != nil {
			p.errors <- errors.Annotate(err, "Unmarshal [resp]")
			continue
		}

		p.mutex.Lock()
		call, ok := p.pending[resp.ID]
		p.mutex.Unlock()

		if ok {
			p.mutex.Lock()
			delete(p.pending, resp.ID)
			p.mutex.Unlock()

			logging.DDumpJSON("ws resp <", resp)

			if resp.Error != nil {
				call.Error = resp.Error
			}

			call.Reply = resp.Result
			call.done()
		} else {
			var subsResp rpcSubscriptionResponse
			if err := ffjson.Unmarshal([]byte(data), &subsResp); err != nil {
				p.errors <- errors.Annotate(err, "Unmarshal [subsResp]")
				continue
			}

			logging.DDumpJSON("ws subscription resp <", subsResp)

			if subsResp.Method != "notice" {
				p.errors <- errors.Errorf(
					"rpc subscription: invalid method %q",
					subsResp.Method,
				)

				continue
			}

			parms := subsResp.Params
			subscriberID := uint64(parms[0].(float64))

			p.mutexSubscribe.Lock()
			fn, ok := p.subscrFns[subscriberID]
			p.mutexSubscribe.Unlock()

			if !ok {
				p.errors <- errors.Errorf(
					"hook for subscriber ID %d is undefined",
					subscriberID,
				)
				continue
			}

			if fn != nil {
				if err := fn(parms[1]); err != nil {
					p.errors <- errors.Annotate(err, "subscribe callback error")
					continue
				}
			}
		}
	}

	// Terminate pending calls
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.shutdown.Set()
	for _, call := range p.pending {
		call.Error = ErrShutdown
		call.done()
	}
}

func (p *wsClient) Subscribe(apiID int, method string, fn SubscribeCallback, args ...interface{}) (*json.RawMessage, error) {
	p.mutexSubscribe.Lock()
	if p.subscribeID == math.MaxUint64 {
		p.subscribeID = 0
	}

	p.subscribeID++
	p.subscrFns[p.subscribeID] = fn
	p.mutexSubscribe.Unlock()

	return p.CallAPI(
		apiID, method,
		append([]interface{}{
			p.subscribeID,
		}, args...)...)
}

func (p *wsClient) OnError(fn ErrorFunc) {
	p.onError = fn
}

func (p *wsClient) CallAPI(apiID int, method string, args ...interface{}) (*json.RawMessage, error) {
	call, err := p.Call("call", []interface{}{
		apiID,
		method,
		args,
	})
	if err != nil {
		return nil, errors.Annotate(err, "Call")
	}

	<-call.Done
	return call.Reply, call.Error
}

func (p *wsClient) Call(method string, args []interface{}) (*RPCCall, error) {
	if !p.IsConnected() {
		return nil, ErrShutdown
	}

	call := &RPCCall{
		Request: rpcRequest{
			Method: method,
			Params: args,
			ID:     p.requestID,
		},
		Done: make(chan *RPCCall, 200),
	}

	p.mutex.Lock()
	if p.requestID == math.MaxUint64 {
		p.requestID = 0
	}

	p.requestID++
	p.pending[call.Request.ID] = call
	p.mutex.Unlock()

	logging.DDumpJSON("ws req >", call.Request)

	if err := p.conn.SetDeadline(time.Now().Add(ReadWriteTimeout)); err != nil {
		return nil, errors.Annotate(err, "SetDeadline")
	}

	if err := websocket.JSON.Send(p.conn, call.Request); err != nil {
		p.mutex.Lock()
		delete(p.pending, call.Request.ID)
		p.mutex.Unlock()

		return nil, errors.Annotate(err, "Send [req]")
	}

	return call, nil
}
