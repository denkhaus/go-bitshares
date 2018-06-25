package api

import (
	"log"
	"sync"

	"github.com/denkhaus/bitshares/client"
	"github.com/denkhaus/bitshares/latency"
	"github.com/juju/errors"
)

type ClientProvider interface {
	OnError(fn client.ErrorFunc)
	Connect() error
	OnNotify(subscriberID int, fn client.NotifyFunc) error
	CallAPI(apiID int, method string, args ...interface{}) (interface{}, error)
	SetDebug(debug bool)
	Close() error
}

type SimpleClientProvider struct {
	client.WebsocketClient
	api BitsharesAPI
}

func NewSimpleClientProvider(endpointURL string, api BitsharesAPI) *SimpleClientProvider {
	wsc := client.NewWebsocketClient(endpointURL)
	sim := SimpleClientProvider{
		api:             api,
		WebsocketClient: wsc,
	}

	return &sim
}

func (p *SimpleClientProvider) CallAPI(apiID int, method string, args ...interface{}) (interface{}, error) {
	if !p.IsConnected() {
		if err := p.api.Connect(); err != nil {
			return nil, errors.Annotate(err, "Connect [api]")
		}
	}

	return p.WebsocketClient.CallAPI(apiID, method, args...)
}

type BestNodeClientProvider struct {
	mu  sync.Mutex
	api BitsharesAPI
	client.WebsocketClient
	tester latency.LatencyTester
}

func NewBestNodeClientProvider(endpointURL string, api BitsharesAPI) (*BestNodeClientProvider, error) {
	tester, err := latency.NewLatencyTester(endpointURL)
	if err != nil {
		return nil, errors.Annotate(err, "NewLatencyTester")
	}

	pr := &BestNodeClientProvider{
		api:             api,
		tester:          tester,
		WebsocketClient: tester.TopNodeClient(),
	}

	tester.OnTopNodeChanged(pr.onTopNodeChanged)
	tester.Start()

	return pr, nil
}

func (p *BestNodeClientProvider) onTopNodeChanged(newEndpoint string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	log.Printf("change top node client -> %s\n", newEndpoint)

	if p.IsConnected() {
		if err := p.WebsocketClient.Close(); err != nil {
			return errors.Annotate(err, "Close [client]")
		}
	}

	p.WebsocketClient = p.tester.TopNodeClient()
	return nil
}

func (p *BestNodeClientProvider) CallAPI(apiID int, method string, args ...interface{}) (interface{}, error) {
	if !p.IsConnected() {
		if err := p.api.Connect(); err != nil {
			return nil, errors.Annotate(err, "Connect [api]")
		}
	}

	p.mu.Lock()
	resp, err := p.WebsocketClient.CallAPI(apiID, method, args...)
	p.mu.Unlock()

	return resp, err
}

func (p *BestNodeClientProvider) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.IsConnected() {
		p.WebsocketClient.Close()
		if err := p.WebsocketClient.Close(); err != nil {
			return errors.Annotate(err, "Close [client]")
		}
	}

	if err := p.tester.Close(); err != nil {
		return errors.Annotate(err, "Close [tester]")
	}

	return nil
}
