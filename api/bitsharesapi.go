package api

import (
	"time"

	"github.com/denkhaus/bitshares/objects"
	"github.com/denkhaus/bitshares/rpc"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

const (
	InvalidApiID         = -1
	AssetsListAll        = -1
	AssetsMaxBatchSize   = 100
	GetCallOrdersLimit   = 100
	GetLimitOrdersLimit  = 100
	GetSettleOrdersLimit = 100
	GetTradeHistoryLimit = 100
)

var (
	EmptyParams = []interface{}{}
)

type BitsharesAPI interface {
	Close() error
	Connect() error
	DatabaseApiID() int
	CryptoApiID() int
	HistoryApiID() int
	BroadcastApiID() int
	SetCredentials(username, password string)
	OnError(func(error))
	OnNotify(subscriberID int, notifyFn func(msg interface{}) error) error
	CallWebsocketAPI(apiID int, method string, args ...interface{}) (interface{}, error)
	SetSubscribeCallback(notifyID int, clearFilter bool) error
	CancelAllSubscriptions() error
	CancelOrder(orderID objects.GrapheneObject, broadcast bool) (*objects.Transaction, error)
	GetBlock(number uint64) (*objects.Block, error)
	GetDynamicGlobalProperties() (*objects.DynamicGlobalProperties, error)
	SubscribeToMarket(notifyID int, base objects.GrapheneObject, quote objects.GrapheneObject) error
	UnsubscribeFromMarket(base objects.GrapheneObject, quote objects.GrapheneObject) error
	GetAccountBalances(account objects.GrapheneObject, assets ...objects.GrapheneObject) ([]objects.AssetAmount, error)
	GetAccountByName(name string) (*objects.Account, error)
	GetAccounts(accountIDs ...objects.GrapheneObject) ([]objects.Account, error)
	GetMarginPositions(accountID objects.GrapheneObject) ([]objects.CallOrder, error)
	GetCallOrders(assetID objects.GrapheneObject, limit int) ([]objects.CallOrder, error)
	GetLimitOrders(base, quote objects.GrapheneObject, limit int) (objects.LimitOrders, error)
	GetObjects(objectIDs ...objects.GrapheneObject) ([]interface{}, error)
	GetSettleOrders(assetID objects.GrapheneObject, limit int) ([]objects.SettleOrder, error)
	//Broadcast(wifKeys []string, feeAsset objects.GrapheneObject, ops ...objects.Operation) (string, error)
	GetTradeHistory(base, quote objects.GrapheneObject, toTime, fromTime time.Time, limit int) ([]objects.MarketTrade, error)
	ListAssets(lowerBoundSymbol string, limit int) ([]objects.Asset, error)
	GetChainID() (string, error)

	//wallet API
	//Transfer(from, to objects.GrapheneObject, amount objects.AssetAmount) (interface{}, error)
	ListAccountBalances(account objects.GrapheneObject) ([]objects.AssetAmount, error)
	WalletLock() error
	WalletUnlock(password string) error
	WalletIsLocked() (bool, error)
	BorrowAsset(account objects.GrapheneObject, amountToBorrow string, symbolToBorrow objects.GrapheneObject, amountOfCollateral string, broadcast bool) (*objects.Transaction, error)
	Buy(account objects.GrapheneObject, base, quote objects.GrapheneObject, rate string, amount string, broadcast bool) (*objects.Transaction, error)
	BuyEx(account objects.GrapheneObject, base, quote objects.GrapheneObject, rate float64, amount float64, broadcast bool) (*objects.Transaction, error)
	Sell(account objects.GrapheneObject, base, quote objects.GrapheneObject, rate string, amount string, broadcast bool) (*objects.Transaction, error)
	SellEx(account objects.GrapheneObject, base, quote objects.GrapheneObject, rate float64, amount float64, broadcast bool) (*objects.Transaction, error)
	SellAsset(account objects.GrapheneObject, amountToSell string, symbolToSell objects.GrapheneObject, minToReceive string, symbolToReceive objects.GrapheneObject, timeout uint32, fillOrKill bool, broadcast bool) (*objects.Transaction, error)
}

type bitsharesAPI struct {
	wsClient       rpc.WebsocketClient
	rpcClient      rpc.RPCClient
	chainConfig    *ChainConfig
	username       string
	password       string
	databaseApiID  int
	historyApiID   int
	cryptoApiID    int
	broadcastApiID int
}

func (p *bitsharesAPI) DatabaseApiID() int {
	return p.databaseApiID
}

func (p *bitsharesAPI) BroadcastApiID() int {
	return p.broadcastApiID
}

func (p *bitsharesAPI) HistoryApiID() int {
	return p.historyApiID
}

func (p *bitsharesAPI) CryptoApiID() int {
	return p.cryptoApiID
}

func (p *bitsharesAPI) getApiID(identifier string) (int, error) {
	resp, err := p.wsClient.CallAPI(1, identifier, EmptyParams)
	if err != nil {
		return InvalidApiID, errors.Annotate(err, identifier)
	}

	//util.Dump(identifier+" in", resp)
	return int(resp.(float64)), nil
}

func (p *bitsharesAPI) login() (bool, error) {
	resp, err := p.wsClient.CallAPI(1, "login", p.username, p.password)
	if err != nil {
		return false, err // errors.Annotate(err, "login")
	}

	//util.Dump("login in", resp)
	return resp.(bool), nil
}

func (p *bitsharesAPI) SetSubscribeCallback(notifyID int, clearFilter bool) error {
	// returns nil if successfull
	_, err := p.wsClient.CallAPI(p.databaseApiID, "set_subscribe_callback", notifyID, clearFilter)
	if err != nil {
		return err
	}

	return nil
}

func (p *bitsharesAPI) SubscribeToMarket(notifyID int, base objects.GrapheneObject, quote objects.GrapheneObject) error {
	// returns nil if successfull
	_, err := p.wsClient.CallAPI(p.databaseApiID, "subscribe_to_market", notifyID, base.Id(), quote.Id())
	if err != nil {
		return err
	}

	return nil
}

func (p *bitsharesAPI) UnsubscribeFromMarket(base objects.GrapheneObject, quote objects.GrapheneObject) error {
	// returns nil if successfull
	_, err := p.wsClient.CallAPI(p.databaseApiID, "unsubscribe_from_market", base.Id(), quote.Id())
	if err != nil {
		return err
	}

	return nil
}

func (p *bitsharesAPI) CancelAllSubscriptions() error {
	// returns nil
	_, err := p.wsClient.CallAPI(p.databaseApiID, "cancel_all_subscriptions", EmptyParams)
	if err != nil {
		return err
	}

	return nil
}

//Broadcast a transaction to the network.
//The transaction will be checked for validity in the local database prior to broadcasting.
//If it fails to apply locally, an error will be thrown and the transaction will not be broadcast.
func (p *bitsharesAPI) BroadcastTransaction(tx *objects.Transaction) (string, error) {

	resp, err := p.wsClient.CallAPI(p.BroadcastApiID(), "broadcast_transaction", tx)
	if err != nil {
		return "", err
	}

	util.Dump("broadcast_transaction <", resp)
	return resp.(string), nil
}

//GetPotentialSignatures will return the set of all public keys that could possibly sign for a given transaction.
//This call can be used by wallets to filter their set of public keys to just the relevant subset prior to calling
//GetRequiredSignatures to get the minimum subset.
func (p *bitsharesAPI) GetPotentialSignatures(tx *objects.Transaction) ([]objects.PublicKey, error) {

	resp, err := p.wsClient.CallAPI(p.DatabaseApiID(), "get_potential_signatures", tx)
	if err != nil {
		return nil, err
	}

	util.Dump("get_potential_signatures <", resp)

	data := resp.([]interface{})
	ret := make([]objects.PublicKey, len(data))

	for idx, acct := range data {
		if err := ffjson.Unmarshal(util.ToBytes(acct), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal PublicKey")
		}
	}

	return ret, nil
}

//GetBlock returns a Block by block number.
func (p *bitsharesAPI) GetBlock(number uint64) (*objects.Block, error) {
	resp, err := p.wsClient.CallAPI(0, "get_block", number)
	if err != nil {
		return nil, err // errors.Annotate(err, "get_block")
	}

	//util.Dump("get_block <", resp)
	ret := objects.Block{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Block")
	}

	return &ret, nil
}

//GetAccountByName returns a Account object by username
func (p *bitsharesAPI) GetAccountByName(name string) (*objects.Account, error) {
	resp, err := p.wsClient.CallAPI(0, "get_account_by_name", name)
	if err != nil {
		return nil, err // errors.Annotate(err, "get_account_by_name")
	}

	//util.Dump("get_account_by_name <", resp)
	ret := objects.Account{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Account")
	}

	return &ret, nil
}

//GetAccounts returns a list of accounts by ID.
func (p *bitsharesAPI) GetAccounts(accounts ...objects.GrapheneObject) ([]objects.Account, error) {

	ids := objects.GrapheneObjects(accounts).ToObjectIDs()
	resp, err := p.wsClient.CallAPI(0, "get_accounts", ids)
	if err != nil {
		return nil, err // errors.Annotate(err, "get_accounts")
	}

	data := resp.([]interface{})
	ret := make([]objects.Account, len(data))

	for idx, acct := range data {
		if err := ffjson.Unmarshal(util.ToBytes(acct), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal Account")
		}
	}

	return ret, nil
}

//GetDynamicGlobalProperties
func (p *bitsharesAPI) GetDynamicGlobalProperties() (*objects.DynamicGlobalProperties, error) {

	resp, err := p.wsClient.CallAPI(0, "get_dynamic_global_properties", EmptyParams)
	if err != nil {
		return nil, err
	}

	//util.Dump("get_dynamic_global_properties <", resp)

	var ret objects.DynamicGlobalProperties
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal DynamicGlobalProperties")
	}

	return &ret, nil
}

//GetAccountBalances retrieves AssetAmount objects by given AccountID
func (p *bitsharesAPI) GetAccountBalances(account objects.GrapheneObject, assets ...objects.GrapheneObject) ([]objects.AssetAmount, error) {

	ids := objects.GrapheneObjects(assets).ToObjectIDs()
	resp, err := p.wsClient.CallAPI(0, "get_account_balances", account.Id(), ids)
	if err != nil {
		return nil, err
	}

	data := resp.([]interface{})
	ret := make([]objects.AssetAmount, len(data))

	for idx, a := range data {
		if err := ffjson.Unmarshal(util.ToBytes(a), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal AssetAmount")
		}
	}

	return ret, nil
}

//ListAssets retrieves assets
//@param lowerBoundSymbol: Lower bound of symbol names to retrieve
//@param limit: Maximum number of assets to fetch, if the constant AssetsListAll is passed, all existing assets will be retrieved.
func (p *bitsharesAPI) ListAssets(lowerBoundSymbol string, limit int) ([]objects.Asset, error) {

	lim := limit
	if limit > AssetsMaxBatchSize || limit == AssetsListAll {
		lim = AssetsMaxBatchSize
	}

	resp, err := p.wsClient.CallAPI(0, "list_assets", lowerBoundSymbol, lim)
	if err != nil {
		return nil, err // errors.Annotate(err, "list_assets")
	}

	data := resp.([]interface{})
	ret := make([]objects.Asset, len(data))

	for idx, a := range data {
		if err := ffjson.Unmarshal(util.ToBytes(a), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal Asset")
		}
	}

	return ret, nil
}

//GetRequiredFees calculates the required fee for each operation in the specified asset type.
//If the asset type does not have a valid core_exchange_rate
func (p *bitsharesAPI) GetRequiredFees(ops objects.Operations, feeAsset objects.GrapheneObject) ([]objects.AssetAmount, error) {
	resp, err := p.wsClient.CallAPI(0, "get_required_fees", ops.Types(), feeAsset.Id())
	if err != nil {
		return nil, err
	}

	//util.Dump("get_required_fees <", resp)

	data := resp.([]interface{})
	ret := make([]objects.AssetAmount, len(data))

	for idx, a := range data {
		if err := ffjson.Unmarshal(util.ToBytes(a), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal AssetAmount")
		}
	}

	return ret, nil
}

//GetLimitOrders returns a slice of LimitOrder objects.
func (p *bitsharesAPI) GetLimitOrders(base, quote objects.GrapheneObject, limit int) (objects.LimitOrders, error) {
	if limit > GetLimitOrdersLimit {
		limit = GetLimitOrdersLimit
	}

	resp, err := p.wsClient.CallAPI(0, "get_limit_orders", base.Id(), quote.Id(), limit)
	if err != nil {
		return nil, err
	}

	//util.Dump("limitorders <", resp)

	data := resp.([]interface{})
	ret := make(objects.LimitOrders, len(data))

	for idx, a := range data {
		if err := ffjson.Unmarshal(util.ToBytes(a), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal LimitOrder")
		}
	}

	return ret, nil
}

//GetSettleOrders returns a slice of SettleOrder objects.
func (p *bitsharesAPI) GetSettleOrders(assetID objects.GrapheneObject, limit int) ([]objects.SettleOrder, error) {
	if limit > GetSettleOrdersLimit {
		limit = GetSettleOrdersLimit
	}

	resp, err := p.wsClient.CallAPI(0, "get_settle_orders", assetID.Id(), limit)
	if err != nil {
		return nil, err // errors.Annotate(err, "get_settle_orders")
	}

	//util.Dump("settleorders in", resp)

	data := resp.([]interface{})
	ret := make([]objects.SettleOrder, len(data))

	for idx, a := range data {
		if err := ffjson.Unmarshal(util.ToBytes(a), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal SettleOrder")
		}
	}

	return ret, nil
}

//GetCallOrders returns a slice of CallOrder objects.
func (p *bitsharesAPI) GetCallOrders(assetID objects.GrapheneObject, limit int) ([]objects.CallOrder, error) {
	if limit > GetCallOrdersLimit {
		limit = GetCallOrdersLimit
	}

	resp, err := p.wsClient.CallAPI(0, "get_call_orders", assetID.Id(), limit)
	if err != nil {
		return nil, err
	}

	//util.Dump("callorders in", resp)

	data := resp.([]interface{})
	ret := make([]objects.CallOrder, len(data))

	for idx, a := range data {
		if err := ffjson.Unmarshal(util.ToBytes(a), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal CallOrder")
		}
	}

	return ret, nil
}

//GetMarginPositions returns a slice of CallOrder objects for the specified account.
func (p *bitsharesAPI) GetMarginPositions(accountID objects.GrapheneObject) ([]objects.CallOrder, error) {
	resp, err := p.wsClient.CallAPI(0, "get_margin_positions", accountID.Id())
	if err != nil {
		return nil, err
	}

	//util.Dump("marginpositions in", resp)

	data := resp.([]interface{})
	ret := make([]objects.CallOrder, len(data))

	for idx, a := range data {
		if err := ffjson.Unmarshal(util.ToBytes(a), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal CallOrder")
		}
	}

	return ret, nil
}

//GetTradeHistory returns MarketTrade object.
func (p *bitsharesAPI) GetTradeHistory(base, quote objects.GrapheneObject, toTime, fromTime time.Time, limit int) ([]objects.MarketTrade, error) {
	if limit > GetTradeHistoryLimit {
		limit = GetTradeHistoryLimit
	}

	resp, err := p.wsClient.CallAPI(0, "get_trade_history", base.Id(), quote.Id(), toTime, fromTime, limit)
	if err != nil {
		return nil, err
	}

	data := resp.([]interface{})
	ret := make([]objects.MarketTrade, len(data))

	for idx, a := range data {
		if err := ffjson.Unmarshal(util.ToBytes(a), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal MarketTrade")
		}
	}

	return ret, nil
}

//GetChainID returns the ID of the chain we are connected to.
func (p *bitsharesAPI) GetChainID() (string, error) {
	resp, err := p.wsClient.CallAPI(p.databaseApiID, "get_chain_id", EmptyParams)
	if err != nil {
		return "", err
	}

	//util.Dump("get_chain_id <", resp)
	return resp.(string), nil
}

//GetObjects returns a list of Graphene Objects by ID.
func (p *bitsharesAPI) GetObjects(ids ...objects.GrapheneObject) ([]interface{}, error) {
	params := objects.GrapheneObjects(ids).ToObjectIDs()
	resp, err := p.wsClient.CallAPI(0, "get_objects", params)
	if err != nil {
		return nil, err
	}

	//	util.Dump("get_objects <", resp)

	data := resp.([]interface{})
	ret := make([]interface{}, len(data))
	id := objects.GrapheneID{}

	for idx, obj := range data {
		if obj == nil {
			continue
		}

		if err := id.FromRawData(obj); err != nil {
			return nil, errors.Annotate(err, "from raw data")
		}

		b := util.ToBytes(obj)

		switch id.Space() {
		case objects.SpaceTypeProtocol:
			switch id.Type() {
			case objects.ObjectTypeAsset:
				ass := objects.Asset{}
				if err := ffjson.Unmarshal(b, &ass); err != nil {
					return nil, errors.Annotate(err, "unmarshal Asset")
				}
				ret[idx] = ass

			case objects.ObjectTypeAccount:
				acc := objects.Account{}
				if err := ffjson.Unmarshal(b, &acc); err != nil {
					return nil, errors.Annotate(err, "unmarshal Account")
				}
				ret[idx] = acc

			case objects.ObjectTypeForceSettlement:
				set := objects.SettleOrder{}
				if err := ffjson.Unmarshal(b, &set); err != nil {
					return nil, errors.Annotate(err, "unmarshal SettleOrder")
				}
				ret[idx] = set

			case objects.ObjectTypeLimitOrder:
				lim := objects.LimitOrder{}
				if err := ffjson.Unmarshal(b, &lim); err != nil {
					return nil, errors.Annotate(err, "unmarshal LimitOrder")
				}
				ret[idx] = lim

			case objects.ObjectTypeCallOrder:
				cal := objects.CallOrder{}
				if err := ffjson.Unmarshal(b, &cal); err != nil {
					return nil, errors.Annotate(err, "unmarshal CallOrder")
				}
				ret[idx] = cal
			default:
				return nil, errors.Errorf("unable to parse GrapheneObject with ID %s", id)
			}
		case objects.SpaceTypeImplementation:
			switch id.Type() {
			case objects.ObjectTypeAssetBitAssetData:
				bit := objects.BitAssetData{}
				if err := ffjson.Unmarshal(b, &bit); err != nil {
					return nil, errors.Annotate(err, "unmarshal BitAssetData")
				}
				ret[idx] = bit

			default:
				return nil, errors.Errorf("unable to parse GrapheneObject with ID %s", id)
			}
		}
	}

	return ret, nil
}

func (p *bitsharesAPI) CancelOrder(orderID objects.GrapheneObject, broadcast bool) (*objects.Transaction, error) {
	resp, err := p.wsClient.CallAPI(0, "cancel_order", orderID.Id(), broadcast)
	if err != nil {
		return nil, err
	}

	ret := objects.Transaction{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Transaction")
	}

	return &ret, nil
}

func (p *bitsharesAPI) CallWebsocketAPI(apiID int, method string, args ...interface{}) (interface{}, error) {
	return p.wsClient.CallAPI(apiID, method, args...)
}

func (p *bitsharesAPI) OnNotify(subscriberID int, notifyFn func(msg interface{}) error) error {
	return p.wsClient.OnNotify(subscriberID, notifyFn)
}

func (p *bitsharesAPI) OnError(errorFn func(err error)) {
	p.wsClient.OnError(errorFn)
}

//SetCredentials defines username and password for login.
func (p *bitsharesAPI) SetCredentials(username, password string) {
	p.username = username
	p.password = password
}

func (p *bitsharesAPI) Connect() (err error) {
	if err := p.wsClient.Connect(); err != nil {
		return errors.Annotate(err, "ws connect")
	}

	if err := p.rpcClient.Connect(); err != nil {
		return errors.Annotate(err, "rpc connect")
	}

	if ok, err := p.login(); err != nil || !ok {
		if err != nil {
			return errors.Annotate(err, "login")
		}
		return errors.New("login not successful")
	}

	p.databaseApiID, err = p.getApiID("database")
	if err != nil {
		return errors.Annotate(err, "get database API ID")
	}

	p.historyApiID, err = p.getApiID("history")
	if err != nil {
		return errors.Annotate(err, "get history API ID")
	}

	p.broadcastApiID, err = p.getApiID("network_broadcast")
	if err != nil {
		return errors.Annotate(err, "get network API ID")
	}

	p.cryptoApiID, err = p.getApiID("crypto")
	if err != nil {
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

	return nil
}

//Close() shuts down the api and closes underlying clients.
func (p *bitsharesAPI) Close() error {

	if p.rpcClient != nil {
		if err := p.rpcClient.Close(); err != nil {
			return errors.Annotate(err, "close rpc client")
		}
		p.rpcClient = nil
	}

	if p.wsClient != nil {
		if err := p.wsClient.Close(); err != nil {
			return errors.Annotate(err, "close ws client")
		}
		p.wsClient = nil
	}

	return nil
}

//New creates a new BitsharesAPI interface.
func New(wsEndpointURL, rpcEndpointURL string) BitsharesAPI {
	api := bitsharesAPI{
		wsClient:       rpc.NewWebsocketClient(wsEndpointURL),
		rpcClient:      rpc.NewRPCClient(rpcEndpointURL),
		databaseApiID:  InvalidApiID,
		historyApiID:   InvalidApiID,
		broadcastApiID: InvalidApiID,
		cryptoApiID:    InvalidApiID,
	}

	return &api
}
