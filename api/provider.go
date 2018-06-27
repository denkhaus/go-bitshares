package api

import (
	"sync"

	"github.com/denkhaus/logging"
	"github.com/juju/errors"
)

type ClientProvider interface {
	OnError(fn ErrorFunc)
	Connect() error
	OnNotify(subscriberID int, fn NotifyFunc) error
	CallAPI(apiID int, method string, args ...interface{}) (interface{}, error)
	Close() error
}

type SimpleClientProvider struct {
	WebsocketClient
	api BitsharesAPI
}

func NewSimpleClientProvider(endpointURL string, api BitsharesAPI) ClientProvider {
	wsc := NewWebsocketClient(endpointURL)
	sim := SimpleClientProvider{
		api:             api,
		WebsocketClient: wsc,
	}

	return &sim
}

func (p *SimpleClientProvider) CallAPI(apiID int, method string, args ...interface{}) (interface{}, error) {
	if !p.WebsocketClient.IsConnected() {
		if err := p.api.Connect(); err != nil {
			return nil, errors.Annotate(err, "Connect [api]")
		}
	}

	return p.WebsocketClient.CallAPI(apiID, method, args...)
}

type BestNodeClientProvider struct {
	WebsocketClient
	mu     sync.Mutex
	api    BitsharesAPI
	tester LatencyTester
}

func NewBestNodeClientProvider(endpointURL string, api BitsharesAPI) (ClientProvider, error) {
	tester, err := NewLatencyTester(endpointURL)
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

	logging.Debugf("change top node client -> %s\n", newEndpoint)

	if p.WebsocketClient.IsConnected() {
		logging.Debug("close [client]")
		if err := p.WebsocketClient.Close(); err != nil {
			return errors.Annotate(err, "Close [client]")
		}
	}

	p.WebsocketClient = p.tester.TopNodeClient()
	return nil
}

func (p *BestNodeClientProvider) CallAPI(apiID int, method string, args ...interface{}) (interface{}, error) {

	p.mu.Lock()
	conn := p.WebsocketClient.IsConnected()
	p.mu.Unlock()

	if !conn {
		logging.Debug("reconnect api")
		if err := p.api.Connect(); err != nil {
			return nil, errors.Annotate(err, "Connect [api]")
		}
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	return p.WebsocketClient.CallAPI(apiID, method, args...)
}

func (p *BestNodeClientProvider) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	logging.Debug("close provider")
	if p.WebsocketClient.IsConnected() {
		logging.Debug("close [client]")
		if err := p.WebsocketClient.Close(); err != nil {
			return errors.Annotate(err, "Close [client]")
		}
	}

	logging.Debug("close [tester]")
	if err := p.tester.Close(); err != nil {
		return errors.Annotate(err, "Close [tester]")
	}

	return nil
}
