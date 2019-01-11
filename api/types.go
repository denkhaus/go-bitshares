package api

import "log"

type NotifyFunc func(msg interface{}) error
type ErrorFunc func(error)

type WebsocketClient interface {
	IsConnected() bool
	OnError(fn ErrorFunc)
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

func (p rpcResponse) HasError() bool {
	return p.Error.Code != 0
}

type rpcNotify struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}
