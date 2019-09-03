package bitshares

import (
	"encoding/json"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/logging"
	"github.com/juju/errors"
	deadlock "github.com/sasha-s/go-deadlock"
	"github.com/tevino/abool"
)

type ClientProvider interface {
	OnError(fn api.ErrorFunc)
	Connect() error
	Subscribe(apiID int, method string, fn api.SubscribeCallback, args ...interface{}) (*json.RawMessage, error)
	CallAPI(apiID int, method string, args ...interface{}) (*json.RawMessage, error)
	Close() error
}

type SimpleClientProvider struct {
	api.WebsocketClient
	api WebsocketAPI
}

func NewSimpleClientProvider(endpointURL string, ws WebsocketAPI) ClientProvider {
	wsc := api.NewWebsocketClient(endpointURL)
	sim := SimpleClientProvider{
		api:             ws,
		WebsocketClient: wsc,
	}

	return &sim
}

func (p *SimpleClientProvider) CallAPI(apiID int, method string, args ...interface{}) (*json.RawMessage, error) {
	if !p.WebsocketClient.IsConnected() {
		if err := p.api.Connect(); err != nil {
			return nil, errors.Annotate(err, "Connect [api]")
		}
	}

	rawMessage, err := p.WebsocketClient.CallAPI(apiID, method, args...)
	if err != nil {
		return nil, errors.Annotate(err, "call api error")
	}

	if rawMessage == nil {
		return nil, errors.New("raw message returned is nil")
	}

	return rawMessage, err
}

type BestNodeClientProvider struct {
	api.WebsocketClient
	mu          deadlock.RWMutex
	nodeChanged *abool.AtomicBool
	api         WebsocketAPI
	tester      api.LatencyTester
}

func NewBestNodeClientProvider(endpointURL string, ws WebsocketAPI) (ClientProvider, error) {
	tester, err := api.NewLatencyTester(endpointURL)
	if err != nil {
		return nil, errors.Annotate(err, "NewLatencyTester")
	}

	pr := &BestNodeClientProvider{
		api:             ws,
		tester:          tester,
		nodeChanged:     abool.NewBool(false),
		WebsocketClient: tester.TopNodeClient(),
	}

	tester.OnTopNodeChanged(pr.onTopNodeChanged)
	tester.Start()

	return pr, nil
}

func (p *BestNodeClientProvider) onTopNodeChanged(newEndpoint string) error {
	logging.Infof("change top node client -> %s\n", newEndpoint)
	p.nodeChanged.Set()
	return nil
}

func (p *BestNodeClientProvider) renewClient() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.WebsocketClient.IsConnected() {
		logging.Debug("close [client]")
		if err := p.WebsocketClient.Close(); err != nil {
			return errors.Annotate(err, "Close [client]")
		}
	}

	p.WebsocketClient = p.tester.TopNodeClient()
	return nil
}

func (p *BestNodeClientProvider) handleReconnect() error {
	if err := p.renewClient(); err != nil {
		return errors.Annotate(err, "renewClient")
	}

	logging.Debug("reconnect api")
	if err := p.api.Connect(); err != nil {
		return errors.Annotate(err, "Connect [api]")
	}

	return nil
}

func (p *BestNodeClientProvider) needsReconnect() bool {
	return p.nodeChanged.IsSet() ||
		!p.WebsocketClient.IsConnected()
}

func (p *BestNodeClientProvider) CallAPI(apiID int, method string, args ...interface{}) (*json.RawMessage, error) {
	//ensure reliable connection
	if p.needsReconnect() {
		// either way unsignal nodeChanged
		p.nodeChanged.UnSet()
		if err := p.handleReconnect(); err != nil {
			return nil, errors.Annotate(err, "handleReconnect")
		}
	}

	p.mu.RLock()
	defer p.mu.RUnlock()
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
