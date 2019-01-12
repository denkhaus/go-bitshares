package api

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/denkhaus/logging"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type RPCClient interface {
	CallAPI(method string, args ...interface{}) (*json.RawMessage, error)
	Close() error
	Connect() error
}

type rpcClient struct {
	*http.Client
	*ffjson.Encoder
	*ffjson.Decoder
	encBuf      *bytes.Buffer
	endpointURL string
	timeout     int
}

func (p *rpcClient) Connect() error {
	p.Client = &http.Client{
		Timeout: 10 * time.Second,
	}

	p.encBuf = new(bytes.Buffer)
	p.Encoder = ffjson.NewEncoder(p.encBuf)
	p.Decoder = ffjson.NewDecoder()

	return nil
}

func (p *rpcClient) Close() error {
	return nil
}

func (p *rpcClient) CallAPI(method string, args ...interface{}) (*json.RawMessage, error) {
	req := rpcRequest{
		Method: method,
		ID:     uint64(rand.Int63()),
		Params: args,
	}

	if err := p.Encode(req); err != nil {
		return nil, errors.Annotate(err, "Encode")
	}

	logging.DDumpJSON("rpc req >", req)

	r, err := http.NewRequest("POST", p.endpointURL, p.encBuf)
	if err != nil {
		return nil, errors.Annotate(err, "NewRequest")
	}

	r.Close = true
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")

	resp, err := p.Do(r)
	if err != nil {
		return nil, errors.Annotate(err, "do request")
	}

	defer resp.Body.Close()

	var ret rpcResponseString
	if err := p.DecodeReader(resp.Body, &ret); err != nil {
		return nil, errors.Annotate(err, "Decode")
	}

	if ret.Error != nil {
		return nil, ret.Error
	}

	logging.DDumpJSON("rpc resp <", ret.Result)
	return ret.Result, nil
}

//NewRPCClient creates a new RPC Client
func NewRPCClient(rpcEndpointURL string) RPCClient {
	cli := rpcClient{
		endpointURL: rpcEndpointURL,
	}

	return &cli
}
