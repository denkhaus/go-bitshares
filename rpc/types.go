package rpc

import "log"

type NotifyFunc func(msg interface{}) error
type ErrorFunc func(error)

type WebsocketClient interface {
	OnError(fn ErrorFunc)
	OnNotify(subscriberID int, fn NotifyFunc) error
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

func (p RPCResponse) Is(in interface{}) bool {
	if data, ok := in.(map[string]interface{}); ok {
		if _, ok := data["id"]; ok {
			if _, ok := data["result"]; ok {
				return true
			}
		}
	}
	return false
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

func (p RPCNotify) Is(in interface{}) bool {
	if data, ok := in.(map[string]interface{}); ok {
		if _, ok := data["method"]; ok {
			if _, ok := data["params"]; ok {
				return true
			}
		}
	}
	return false
}

func (p *RPCNotify) reset() {
	p.Method = ""
	p.Params = nil
}
