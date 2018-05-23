package client

import "log"

type NotifyFunc func(msg interface{}) error
type ErrorFunc func(error)

type WebsocketClient interface {
	OnError(fn ErrorFunc)
	SetDebug(debug bool)
	OnNotify(subscriberID int, fn NotifyFunc) error
	Call(method string, args []interface{}) (*RPCCall, error)
	CallAPI(apiID int, method string, args ...interface{}) (interface{}, error)
	Close() error
	Connect() error
}

type RPCCall struct {
	Method  string
	Request rpcRequest
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

type rpcRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	ID     uint64        `json:"id"`
}

func (p *rpcRequest) reset() {
	p.ID = 0
	p.Params = nil
	p.Method = ""
}

type ResponseErrorContext struct {
	Level      string `json:"level"`
	File       string `json:"file"`
	Line       int    `json:"line"`
	Method     string `json:"method"`
	Hostname   string `json:"hostname"`
	ThreadName string `json:"thread_name"`
	Timestamp  string `json:"timestamp"`
}
type ResponseErrorStack struct {
	Context ResponseErrorContext `json:"context"`
	Format  string               `json:"format"`
	Data    interface{}          `json:"data"`
}

type ResponseErrorData struct {
	Code    int                  `json:"code"`
	Name    string               `json:"name"`
	Message string               `json:"message"`
	Stack   []ResponseErrorStack `json:"stack"`
}

type ResponseError struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    ResponseErrorData `json:"data"`
}

func (p ResponseError) Error() string {
	return p.Message
}

//wallet API uses string id ???
type rpcResponseString struct {
	ID     string        `json:"id"`
	Result interface{}   `json:"result,omitempty"`
	Error  ResponseError `json:"error"`
}

func (p rpcResponseString) HasError() bool {
	return p.Error.Code != 0
}

type rpcResponse struct {
	ID     uint64        `json:"id"`
	Result interface{}   `json:"result"`
	Error  ResponseError `json:"error"`
}

func (p rpcResponse) Is(in interface{}) bool {
	if data, ok := in.(map[string]interface{}); ok {
		if _, ok := data["id"]; ok {
			if _, ok := data["result"]; ok {
				return true
			}
			if _, ok := data["error"]; ok {
				return true
			}
		}
	}
	return false
}

func (p rpcResponse) HasError() bool {
	return p.Error.Code != 0
}

func (p *rpcResponse) reset() {
	p.Error = ResponseError{}
	p.Result = nil
	p.ID = 0
}

type rpcNotify struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

func (p rpcNotify) Is(in interface{}) bool {
	if data, ok := in.(map[string]interface{}); ok {
		if _, ok := data["method"]; ok {
			if _, ok := data["params"]; ok {
				return true
			}
		}
	}
	return false
}

func (p *rpcNotify) reset() {
	p.Method = ""
	p.Params = nil
}
