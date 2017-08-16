package rpc

import "log"

type WebsocketClient interface {
	OnError(fn func(error))
	OnNotify(subscriberID int, notifyFn func(msg interface{}) error) error
	Call(method string, args interface{}) (*RPCCall, error)
	CallAPI(apiID int, method string, args ...interface{}) (interface{}, error)
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

func (p *RPCResponse) reset() {
	p.ID = 0
	p.Error = nil
	p.Result = nil
}

type RPCNotify struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

func (p *RPCNotify) reset() {
	p.Method = ""
	p.Params = nil
}
