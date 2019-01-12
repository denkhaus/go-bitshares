package api

//go:generate ffjson $GOFILE

import (
	"encoding/json"
	"log"
)

type SubscribeCallback func(msg interface{}) error
type ErrorFunc func(error)

type WebsocketClient interface {
	IsConnected() bool
	OnError(fn ErrorFunc)
	OnSubscribe(subscriberID uint64, fn SubscribeCallback) error
	Call(method string, args []interface{}) (*RPCCall, error)
	CallAPI(apiID int, method string, args ...interface{}) (*json.RawMessage, error)
	Close() error
	Connect() error
}
type RPCCall struct {
	Method  string
	Request rpcRequest
	Reply   *json.RawMessage
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
	ID     string           `json:"id"`
	Result *json.RawMessage `json:"result,omitempty"`
	Error  *ResponseError   `json:"error,omitempty"`
}

type rpcResponse struct {
	ID     uint64           `json:"id"`
	Result *json.RawMessage `json:"result,omitempty"`
	Error  *ResponseError   `json:"error,omitempty"`
}

type rpcSubscriptionResponse struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}
