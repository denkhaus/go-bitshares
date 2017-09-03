package rpc

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/juju/errors"
)

type RPCClient interface {
	CallAPI(method string, args ...interface{}) (interface{}, error)
	Close() error
	Connect() error
}

type rpcClient struct {
	*http.Client
	*json.Encoder

	decBuf      *bytes.Buffer
	endpointURL string
	req         rpcRequest
	res         rpcResponseString
	timeout     int
}

func (p *rpcClient) Connect() error {
	p.Client = &http.Client{
		Timeout: 5 * time.Second,
	}

	p.decBuf = new(bytes.Buffer)
	p.Encoder = json.NewEncoder(p.decBuf)

	return nil
}

func (p *rpcClient) Close() error {
	return nil
}

func (p *rpcClient) CallAPI(method string, args ...interface{}) (interface{}, error) {
	//this "ratelimiter" helps to avoid failed calls. Don't know why api is so slow
	time.Sleep(500 * time.Millisecond)

	p.req.Method = method
	p.req.ID = uint64(rand.Int63())
	p.req.Params = args

	if err := p.Encode(&p.req); err != nil {
		return nil, errors.Annotate(err, "encode")
	}

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

	if err := json.NewDecoder(resp.Body).Decode(&p.res); err != nil {
		return nil, errors.Annotate(err, "Decode")
	}

	if p.res.Error != nil {
		return p.res.Result, errors.Errorf("%v", p.res.Error)
	}

	//util.Dump("rpc resp", resp.Body)

	return p.res.Result, nil
}

//NewRPCClient creates a RPC Client
func NewRPCClient(rpcEndpointURL string) RPCClient {
	cli := rpcClient{
		endpointURL: rpcEndpointURL,
	}

	return &cli
}
