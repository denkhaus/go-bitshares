package api

import (
	"bytes"
	"math/rand"
	"net/http"
	"time"

	"github.com/denkhaus/logging"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type RPCClient interface {
	CallAPI(method string, args ...interface{}) (interface{}, error)
	Close() error
	Connect() error
}

type rpcClient struct {
	*http.Client
	*ffjson.Encoder
	*ffjson.Decoder

	decBuf      *bytes.Buffer
	endpointURL string
	timeout     int
}

func (p *rpcClient) Connect() error {
	p.Client = &http.Client{
		Timeout: 10 * time.Second,
	}

	p.decBuf = new(bytes.Buffer)
	p.Encoder = ffjson.NewEncoder(p.decBuf)
	p.Decoder = ffjson.NewDecoder()

	return nil
}

func (p *rpcClient) Close() error {
	return nil
}

func (p *rpcClient) CallAPI(method string, args ...interface{}) (interface{}, error) {
	request := rpcRequest{
		ID:     uint64(rand.Int63()),
		Method: method,
		Params: args,
	}

	if err := p.Encode(&request); err != nil {
		return nil, errors.Annotate(err, "Encode")
	}

	logging.DDumpJSON("rpc req >", request)

	req, err := http.NewRequest("POST", p.endpointURL, p.decBuf)
	if err != nil {
		return nil, errors.Annotate(err, "NewRequest")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := p.Do(req)
	if err != nil {
		return nil, errors.Annotate(err, "do request")
	}

	defer resp.Body.Close()
	var res rpcResponseString
	if err := p.DecodeReader(resp.Body, &res); err != nil {
		return nil, errors.Annotate(err, "Decode")
	}

	if res.HasError() {
		return res.Result, res.Error
	}

	logging.DDumpJSON("rpc resp <", res.Result)
	return res.Result, nil
}

//NewRPCClient creates a new RPC Client
func NewRPCClient(rpcEndpointURL string) RPCClient {
	cli := rpcClient{
		endpointURL: rpcEndpointURL,
	}

	return &cli
}
