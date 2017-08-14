package api

import (
	"github.com/denkhaus/bitshares/client"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

const (
	InvalidApiID = -1
)

var (
	EmptyParams = []interface{}{}
)

type BitsharesApi struct {
	client        *client.WSClient
<<<<<<< HEAD
	chainConfig   *ChainConfig
=======
>>>>>>> 90972d81a2199b7398b4ac4858bba2c236601463
	username      string
	password      string
	databaseApiID int
	historyApiID  int
	cryptoApiID   int
	networkApiID  int
}

func (p *BitsharesApi) getApiID(identifier string) (int, error) {
	resp, err := p.client.CallApi(1, identifier, EmptyParams)
	if err != nil {
		return InvalidApiID, errors.Annotate(err, identifier)
	}

<<<<<<< HEAD
	//util.Dump(identifier+" in", resp)
=======
	util.Dump(identifier+" in", resp)
>>>>>>> 90972d81a2199b7398b4ac4858bba2c236601463
	return int(resp.(float64)), nil
}

func (p *BitsharesApi) login() (bool, error) {
	resp, err := p.client.CallApi(1, "login", p.username, p.password)
	if err != nil {
		return false, errors.Annotate(err, "login")
	}

	//util.Dump("login in", resp)
	return resp.(bool), nil
}

func (p *BitsharesApi) ensureInitialized() error {
	if !p.isInitialized() {
		return p.initialize()
	}

	return nil
}

func (p *BitsharesApi) isInitialized() bool {
	if p.client != nil {
		if p.databaseApiID != InvalidApiID &&
			p.networkApiID != InvalidApiID &&
			p.cryptoApiID != InvalidApiID &&
			p.historyApiID != InvalidApiID {
			return true
		}
	}

	return false
}

func (p *BitsharesApi) initialize() (err error) {
	if ok, err := p.login(); err != nil || !ok {
		if err != nil {
			return errors.Annotate(err, "login")
		}
		return errors.New("login not successful")
	}

	p.databaseApiID, err = p.getApiID("database")
	if err != nil {
<<<<<<< HEAD
		return errors.Annotate(err, "get database API ID")
=======
		return errors.Annotate(err, "get database API Id")
>>>>>>> 90972d81a2199b7398b4ac4858bba2c236601463
	}

	p.historyApiID, err = p.getApiID("history")
	if err != nil {
<<<<<<< HEAD
		return errors.Annotate(err, "get history API ID")
=======
		return errors.Annotate(err, "get history API Id")
>>>>>>> 90972d81a2199b7398b4ac4858bba2c236601463
	}

	p.networkApiID, err = p.getApiID("network_broadcast")
	if err != nil {
<<<<<<< HEAD
		return errors.Annotate(err, "get network API ID")
=======
		return errors.Annotate(err, "get network API Id")
>>>>>>> 90972d81a2199b7398b4ac4858bba2c236601463
	}

	p.cryptoApiID, err = p.getApiID("crypto")
	if err != nil {
<<<<<<< HEAD
		return errors.Annotate(err, "get crypto API ID")
	}

	chainID, err := p.GetChainID()
	if err != nil {
		return errors.Annotate(err, "get chain ID")
	}

	p.chainConfig, err = p.GetChainConfig(chainID)
	if err != nil {
		return errors.Annotate(err, "get chain config")
	}

	//util.Dump("chain config", p.chainConfig)
=======
		return errors.Annotate(err, "get network API Id")
	}

>>>>>>> 90972d81a2199b7398b4ac4858bba2c236601463
	return nil
}

func (p *BitsharesApi) SetSubscribeCallback(identifier int, clearFilter bool) error {
	if err := p.ensureInitialized(); err != nil {
		return errors.Annotate(err, "ensure initialized")
	}

	resp, err := p.client.CallApi(p.databaseApiID, "set_subscribe_callback", identifier, clearFilter)
	if err != nil {
		return errors.Annotate(err, "set_subscribe_callback")
	}

	util.Dump("set_subscribe_callback in", resp)
	return nil
}

func (p *BitsharesApi) SetCredentials(username, password string) {
	p.username = username
	p.password = password
}

func (p *BitsharesApi) Close() {
	if p.client != nil {
		p.client.Close()
		p.client = nil
	}
}

func New(url string) (*BitsharesApi, error) {
	client, err := client.NewWebsocketClient(url)
	if err != nil {
		return nil, errors.Annotate(err, "new websocket client")
	}

	api := BitsharesApi{
		client:        client,
		databaseApiID: InvalidApiID,
		historyApiID:  InvalidApiID,
		networkApiID:  InvalidApiID,
		cryptoApiID:   InvalidApiID,
	}

	return &api, nil
}
