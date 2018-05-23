package client

import (
	"bytes"
	"math/rand"
	"net/http"
	"time"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type RPCClient interface {
	CallAPI(method string, args ...interface{}) (interface{}, error)
	SetDebug(debug bool)
	Close() error
	Connect() error
}

type rpcClient struct {
	*http.Client
	*ffjson.Encoder
	*ffjson.Decoder

	decBuf      *bytes.Buffer
	endpointURL string
	req         rpcRequest
	res         rpcResponseString
	debug       bool
	timeout     int
}

func (p *rpcClient) Connect() error {
	p.Client = &http.Client{
		Timeout: 3 * time.Second,
	}

	p.decBuf = new(bytes.Buffer)
	p.Encoder = ffjson.NewEncoder(p.decBuf)
	p.Decoder = ffjson.NewDecoder()

	return nil
}

func (p *rpcClient) Close() error {
	return nil
}

func (p *rpcClient) SetDebug(debug bool) {
	p.debug = debug
}

func (p rpcClient) Debug(descr string, in interface{}) {
	if p.debug {
		util.DumpJSON(descr, in)
	}
}

func (p *rpcClient) CallAPI(method string, args ...interface{}) (interface{}, error) {
	p.req.Method = method
	p.req.ID = uint64(rand.Int63())
	p.req.Params = args

	if err := p.Encode(&p.req); err != nil {
		return nil, errors.Annotate(err, "Encode")
	}

	p.Debug("rpc req >", p.req)

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

	if err := p.DecodeReader(resp.Body, &p.res); err != nil {
		return nil, errors.Annotate(err, "Decode")
	}

	if p.res.HasError() {
		return p.res.Result, p.res.Error
	}

	p.Debug("rpc resp <", p.res.Result)

	return p.res.Result, nil
}

//NewRPCClient creates a new RPC Client
func NewRPCClient(rpcEndpointURL string) RPCClient {
	cli := rpcClient{
		endpointURL: rpcEndpointURL,
		debug:       false,
	}

	return &cli
}
