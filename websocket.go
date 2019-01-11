package bitshares

import (
	"time"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/config"
	"github.com/denkhaus/bitshares/crypto"
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/denkhaus/logging"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
	// init operations
	_ "github.com/denkhaus/bitshares/operations"
)

const (
	InvalidApiID                  = -1
	AssetsListAll                 = -1
	AssetsMaxBatchSize            = 100
	GetCallOrdersLimit            = 100
	GetLimitOrdersLimit           = 100
	GetForceSettlementOrdersLimit = 100
	GetTradeHistoryLimit          = 100
	GetAccountHistoryLimit        = 100
)

type WebsocketAPI interface {
	//Common functions
	CallWsAPI(apiID int, method string, args ...interface{}) (interface{}, error)
	Close() error
	Connect() error
	DatabaseAPIID() int
	HistoryAPIID() int
	BroadcastAPIID() int
	SetCredentials(username, password string)
	OnError(api.ErrorFunc)
	OnNotify(subscriberID int, notifyFn func(msg interface{}) error) error
	BuildSignedTransaction(keyBag *crypto.KeyBag, feeAsset types.GrapheneObject, ops ...types.Operation) (*types.SignedTransaction, error)
	SignTransaction(keyBag *crypto.KeyBag, trx *types.SignedTransaction) error

	//Websocket API functions
	BroadcastTransaction(tx *types.SignedTransaction) error
	CancelAllSubscriptions() error
	CancelOrder(orderID types.GrapheneObject, broadcast bool) (*types.SignedTransaction, error)
	GetAccountBalances(account types.GrapheneObject, assets ...types.GrapheneObject) (types.AssetAmounts, error)
	GetAccountByName(name string) (*types.Account, error)
	GetAccountHistory(account types.GrapheneObject, stop types.GrapheneObject, limit int, start types.GrapheneObject) (types.OperationHistories, error)
	GetAccounts(accountIDs ...types.GrapheneObject) (types.Accounts, error)
	GetFullAccounts(accountIDs ...types.GrapheneObject) (types.FullAccountInfos, error)
	GetBlock(number uint64) (*types.Block, error)
	GetCallOrders(assetID types.GrapheneObject, limit int) (types.CallOrders, error)
	GetChainID() (string, error)
	GetDynamicGlobalProperties() (*types.DynamicGlobalProperties, error)
	GetLimitOrders(base, quote types.GrapheneObject, limit int) (types.LimitOrders, error)
	GetOrderBook(base, quote types.GrapheneObject, depth int) (types.OrderBook, error)
	GetMarginPositions(accountID types.GrapheneObject) (types.CallOrders, error)
	GetObjects(objectIDs ...types.GrapheneObject) ([]interface{}, error)
	GetPotentialSignatures(tx *types.SignedTransaction) (types.PublicKeys, error)
	GetRequiredSignatures(tx *types.SignedTransaction, keys types.PublicKeys) (types.PublicKeys, error)
	GetRequiredFees(ops types.Operations, feeAsset types.GrapheneObject) (types.AssetAmounts, error)
	GetForceSettlementOrders(assetID types.GrapheneObject, limit int) (types.ForceSettlementOrders, error)
	GetTradeHistory(base, quote types.GrapheneObject, toTime, fromTime time.Time, limit int) (types.MarketTrades, error)
	ListAssets(lowerBoundSymbol string, limit int) (types.Assets, error)
	SetSubscribeCallback(notifyID int, clearFilter bool) error
	SubscribeToMarket(notifyID int, base types.GrapheneObject, quote types.GrapheneObject) error
	UnsubscribeFromMarket(base types.GrapheneObject, quote types.GrapheneObject) error
	Get24Volume(base types.GrapheneObject, quote types.GrapheneObject) (types.Volume24, error)
}

type websocketAPI struct {
	wsClient       ClientProvider
	username       string
	password       string
	databaseAPIID  int
	historyAPIID   int
	broadcastAPIID int
}

func (p *websocketAPI) getAPIID(identifier string) (int, error) {
	resp, err := p.wsClient.CallAPI(1, identifier, types.EmptyParams)
	if err != nil {
		return InvalidApiID, errors.Annotate(err, identifier)
	}

	logging.DDumpJSON("getApiID <", resp)

	return int(resp.(float64)), nil
}

// login
func (p *websocketAPI) login() (bool, error) {
	resp, err := p.wsClient.CallAPI(1, "login", p.username, p.password)
	if err != nil {
		return false, err
	}

	logging.DDumpJSON("login <", resp)

	return resp.(bool), nil
}

// SetSubscribeCallback
func (p *websocketAPI) SetSubscribeCallback(notifyID int, clearFilter bool) error {
	// returns nil if successful
	_, err := p.wsClient.CallAPI(p.databaseAPIID, "set_subscribe_callback", notifyID, clearFilter)
	if err != nil {
		return err
	}

	return nil
}

// SubscribeToMarket
func (p *websocketAPI) SubscribeToMarket(notifyID int, base types.GrapheneObject, quote types.GrapheneObject) error {
	// returns nil if successful
	_, err := p.wsClient.CallAPI(p.databaseAPIID, "subscribe_to_market", notifyID, base.ID(), quote.ID())
	if err != nil {
		return err
	}

	return nil
}

// UnsubscribeFromMarket
func (p *websocketAPI) UnsubscribeFromMarket(base types.GrapheneObject, quote types.GrapheneObject) error {
	// returns nil if successful
	_, err := p.wsClient.CallAPI(p.databaseAPIID, "unsubscribe_from_market", base.ID(), quote.ID())
	if err != nil {
		return err
	}

	return nil
}

// CancelAllSubscriptions
func (p *websocketAPI) CancelAllSubscriptions() error {
	// returns nil
	_, err := p.wsClient.CallAPI(p.databaseAPIID, "cancel_all_subscriptions", types.EmptyParams)
	if err != nil {
		return err
	}

	return nil
}

//Broadcast a transaction to the network.
//The transaction will be checked for validity prior to broadcasting.
//If it fails to apply at the connected node, an error will be thrown and the transaction will not be broadcast.
func (p *websocketAPI) BroadcastTransaction(tx *types.SignedTransaction) error {
	_, err := p.wsClient.CallAPI(p.broadcastAPIID, "broadcast_transaction", tx)
	if err != nil {
		return err
	}

	return nil
}

//SignTransaction signs a given transaction.
//Required signing keys get selected by API and have to be in keyBag.
func (p *websocketAPI) SignTransaction(keyBag *crypto.KeyBag, tx *types.SignedTransaction) error {
	reqPk, err := p.RequiredSigningKeys(tx)
	if err != nil {
		return errors.Annotate(err, "RequiredSigningKeys")
	}

	signer := crypto.NewTransactionSigner(tx)

	privKeys := keyBag.PrivatesByPublics(reqPk)
	if len(privKeys) == 0 {
		return types.ErrNoSigningKeyFound
	}

	if err := signer.Sign(privKeys, config.Current()); err != nil {
		return errors.Annotate(err, "Sign")
	}

	return nil
}

//BuildSignedTransaction builds a new transaction by given operation(s),
//applies fees, current block data and signs the transaction.
func (p *websocketAPI) BuildSignedTransaction(keyBag *crypto.KeyBag, feeAsset types.GrapheneObject, ops ...types.Operation) (*types.SignedTransaction, error) {
	operations := types.Operations(ops)
	fees, err := p.GetRequiredFees(operations, feeAsset)
	if err != nil {
		return nil, errors.Annotate(err, "GetRequiredFees")
	}

	if err := operations.ApplyFees(fees); err != nil {
		return nil, errors.Annotate(err, "ApplyFees")
	}

	props, err := p.GetDynamicGlobalProperties()
	if err != nil {
		return nil, errors.Annotate(err, "GetDynamicGlobalProperties")
	}

	tx, err := types.NewSignedTransactionWithBlockData(props)
	if err != nil {
		return nil, errors.Annotate(err, "NewTransaction")
	}

	tx.Operations = operations

	reqPk, err := p.RequiredSigningKeys(tx)
	if err != nil {
		return nil, errors.Annotate(err, "RequiredSigningKeys")
	}

	signer := crypto.NewTransactionSigner(tx)

	privKeys := keyBag.PrivatesByPublics(reqPk)
	if len(privKeys) == 0 {
		return nil, types.ErrNoSigningKeyFound
	}

	if err := signer.Sign(privKeys, config.Current()); err != nil {
		return nil, errors.Annotate(err, "Sign")
	}

	return tx, nil
}

//RequiredSigningKeys is a convenience call to retrieve the minimum subset of public keys to sign a transaction.
//If the transaction is already signed, the result is empty.
func (p *websocketAPI) RequiredSigningKeys(tx *types.SignedTransaction) (types.PublicKeys, error) {
	potPk, err := p.GetPotentialSignatures(tx)
	if err != nil {
		return nil, errors.Annotate(err, "GetPotentialSignatures")
	}

	logging.DDumpJSON("potential pubkeys <", potPk)

	reqPk, err := p.GetRequiredSignatures(tx, potPk)
	if err != nil {
		return nil, errors.Annotate(err, "GetRequiredSignatures")
	}

	logging.DDumpJSON("required pubkeys <", reqPk)

	return reqPk, nil
}

//GetPotentialSignatures will return the set of all public keys that could possibly sign for a given transaction.
//This call can be used by wallets to filter their set of public keys to just the relevant subset prior to calling
//GetRequiredSignatures to get the minimum subset.
func (p *websocketAPI) GetPotentialSignatures(tx *types.SignedTransaction) (types.PublicKeys, error) {
	resp, err := p.wsClient.CallAPI(p.databaseAPIID, "get_potential_signatures", tx)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_potential_signatures <", resp)

	ret := types.PublicKeys{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal PublicKeys")
	}

	return ret, nil
}

//GetRequiredSignatures returns the minimum subset of public keys to sign a transaction.
func (p *websocketAPI) GetRequiredSignatures(tx *types.SignedTransaction, potKeys types.PublicKeys) (types.PublicKeys, error) {
	resp, err := p.wsClient.CallAPI(p.databaseAPIID, "get_required_signatures", tx, potKeys)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_required_signatures <", resp)

	ret := types.PublicKeys{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal PublicKeys")
	}

	return ret, nil
}

//GetBlock returns a Block by block number.
func (p *websocketAPI) GetBlock(number uint64) (*types.Block, error) {
	resp, err := p.wsClient.CallAPI(0, "get_block", number)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_block <", resp)

	ret := types.Block{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Block")
	}

	return &ret, nil
}

//GetAccountByName returns a Account object by username
func (p *websocketAPI) GetAccountByName(name string) (*types.Account, error) {
	resp, err := p.wsClient.CallAPI(0, "get_account_by_name", name)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_account_by_name <", resp)

	ret := types.Account{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Account")
	}

	return &ret, nil
}

// GetAccountHistory returns OperationHistory object(s).
// account: The account whose history should be queried
// stop: ID of the earliest operation to retrieve
// limit: Maximum number of operations to retrieve (must not exceed 100)
// start: ID of the most recent operation to retrieve
func (p *websocketAPI) GetAccountHistory(account types.GrapheneObject, stop types.GrapheneObject, limit int, start types.GrapheneObject) (types.OperationHistories, error) {
	if limit > GetAccountHistoryLimit {
		limit = GetAccountHistoryLimit
	}

	resp, err := p.wsClient.CallAPI(p.historyAPIID, "get_account_history", account.ID(), stop.ID(), limit, start.ID())
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_account_history <", resp)

	ret := types.OperationHistories{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Histories")
	}

	return ret, nil
}

//GetAccounts returns a list of accounts by accountID(s).
func (p *websocketAPI) GetAccounts(accounts ...types.GrapheneObject) (types.Accounts, error) {
	ids := types.GrapheneObjects(accounts).ToStrings()
	resp, err := p.wsClient.CallAPI(0, "get_accounts", ids)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_accounts <", resp)

	ret := types.Accounts{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Accounts")
	}

	return ret, nil
}

//GetDynamicGlobalProperties returns essential runtime properties of bitshares network
func (p *websocketAPI) GetDynamicGlobalProperties() (*types.DynamicGlobalProperties, error) {
	resp, err := p.wsClient.CallAPI(0, "get_dynamic_global_properties", types.EmptyParams)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_dynamic_global_properties <", resp)

	var ret types.DynamicGlobalProperties
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal DynamicGlobalProperties")
	}

	return &ret, nil
}

//GetAccountBalances retrieves AssetAmounts by given AccountID
func (p *websocketAPI) GetAccountBalances(account types.GrapheneObject, assets ...types.GrapheneObject) (types.AssetAmounts, error) {
	ids := types.GrapheneObjects(assets).ToStrings()
	resp, err := p.wsClient.CallAPI(0, "get_account_balances", account.ID(), ids)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_account_balances <", resp)

	ret := types.AssetAmounts{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal AssetAmounts")
	}

	return ret, nil
}

//GetFullAccounts retrieves full account information by given AccountIDs
func (p *websocketAPI) GetFullAccounts(accounts ...types.GrapheneObject) (types.FullAccountInfos, error) {
	ids := types.GrapheneObjects(accounts).ToStrings()
	resp, err := p.wsClient.CallAPI(0, "get_full_accounts", ids, false) //do not subscribe for now
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_full_accounts <", resp)

	ret := types.FullAccountInfos{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal FullAccountInfos: "+string(util.ToBytes(resp)))
	}

	return ret, nil
}

// Get24Volume returns the base:quote assets 24h volume
func (p *websocketAPI) Get24Volume(base, quote types.GrapheneObject) (ret types.Volume24, err error) {
	resp, err := p.wsClient.CallAPI(p.databaseAPIID, "get_24_volume", base.ID(), quote.ID())
	if err != nil {
		return
	}

	logging.DDumpJSON("get_24_volume <", resp)

	if err = ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		err = errors.Annotate(err, "unmarshal Volume24")
		return
	}

	return
}

//ListAssets retrieves assets
//lowerBoundSymbol: Lower bound of symbol names to retrieve
//limit: Maximum number of assets to fetch, if the constant AssetsListAll is passed, all existing assets will be retrieved.
func (p *websocketAPI) ListAssets(lowerBoundSymbol string, limit int) (types.Assets, error) {
	if limit > AssetsMaxBatchSize || limit == AssetsListAll {
		limit = AssetsMaxBatchSize
	}

	resp, err := p.wsClient.CallAPI(0, "list_assets", lowerBoundSymbol, limit)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("list_assets <", resp)

	ret := types.Assets{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Assets")
	}

	return ret, nil
}

//GetRequiredFees calculates the required fee for each operation by the specified asset type.
func (p *websocketAPI) GetRequiredFees(ops types.Operations, feeAsset types.GrapheneObject) (types.AssetAmounts, error) {
	resp, err := p.wsClient.CallAPI(0, "get_required_fees", ops.Envelopes(), feeAsset.ID())
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_required_fees <", resp)

	ret := types.AssetAmounts{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal AssetAmounts")
	}

	return ret, nil
}

//GetLimitOrders returns LimitOrders type.
func (p *websocketAPI) GetLimitOrders(base, quote types.GrapheneObject, limit int) (types.LimitOrders, error) {
	if limit > GetLimitOrdersLimit {
		limit = GetLimitOrdersLimit
	}

	resp, err := p.wsClient.CallAPI(0, "get_limit_orders", base.ID(), quote.ID(), limit)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_limit_orders <", resp)

	ret := types.LimitOrders{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal LimitOrders")
	}

	return ret, nil
}

//GetOrderBook returns the OrderBook for the market base:quote.
func (p *websocketAPI) GetOrderBook(base, quote types.GrapheneObject, depth int) (ret types.OrderBook, err error) {

	resp, err := p.wsClient.CallAPI(0, "get_order_book", base.ID(), quote.ID(), depth)
	if err != nil {
		return
	}

	logging.DDumpJSON("get_order_book <", resp)

	if err = ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		err = errors.Annotate(err, "unmarshal LimitOrders")
		return
	}

	return
}

//GetForceSettlementOrders returns ForceSettlementOrders type.
func (p *websocketAPI) GetForceSettlementOrders(assetID types.GrapheneObject, limit int) (types.ForceSettlementOrders, error) {
	if limit > GetForceSettlementOrdersLimit {
		limit = GetForceSettlementOrdersLimit
	}

	resp, err := p.wsClient.CallAPI(0, "get_settle_orders", assetID.ID(), limit)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_settle_orders <", resp)

	ret := types.ForceSettlementOrders{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal ForceSettlementOrders")
	}

	return ret, nil
}

//GetCallOrders returns CallOrders type.
func (p *websocketAPI) GetCallOrders(assetID types.GrapheneObject, limit int) (types.CallOrders, error) {
	if limit > GetCallOrdersLimit {
		limit = GetCallOrdersLimit
	}

	resp, err := p.wsClient.CallAPI(0, "get_call_orders", assetID.ID(), limit)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_call_orders <", resp)

	ret := types.CallOrders{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal CallOrders")
	}

	return ret, nil
}

//GetMarginPositions returns CallOrders type.
func (p *websocketAPI) GetMarginPositions(accountID types.GrapheneObject) (types.CallOrders, error) {
	resp, err := p.wsClient.CallAPI(0, "get_margin_positions", accountID.ID())
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_margin_positions <", resp)

	ret := types.CallOrders{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal CallOrders")
	}

	return ret, nil
}

//GetTradeHistory returns MarketTrades type.
func (p *websocketAPI) GetTradeHistory(base, quote types.GrapheneObject, toTime, fromTime time.Time, limit int) (types.MarketTrades, error) {
	if limit > GetTradeHistoryLimit {
		limit = GetTradeHistoryLimit
	}

	resp, err := p.wsClient.CallAPI(0, "get_trade_history", base.ID(), quote.ID(), toTime, fromTime, limit)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_trade_history <", resp)

	ret := types.MarketTrades{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal MarketTrades")
	}

	return ret, nil
}

//GetChainID returns the ID of the chain we are connected to.
func (p *websocketAPI) GetChainID() (string, error) {
	resp, err := p.wsClient.CallAPI(p.databaseAPIID, "get_chain_id", types.EmptyParams)
	if err != nil {
		return "", err
	}

	logging.DDumpJSON("get_chain_id <", resp)
	return resp.(string), nil
}

//GetObjects returns a list of Graphene Objects by ID.
func (p *websocketAPI) GetObjects(ids ...types.GrapheneObject) ([]interface{}, error) {
	params := types.GrapheneObjects(ids).ToStrings()
	resp, err := p.wsClient.CallAPI(0, "get_objects", params)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_objects <", resp)

	data := resp.([]interface{})
	ret := make([]interface{}, 0)
	id := types.ObjectID{}

	for _, obj := range data {
		if obj == nil {
			continue
		}

		if err := id.FromRawData(obj); err != nil {
			return nil, errors.Annotate(err, "from raw data")
		}

		b := util.ToBytes(obj)

		//TODO: implement
		// ObjectTypeBase
		// ObjectTypeWitness
		// ObjectTypeCustom
		// ObjectTypeProposal
		// ObjectTypeWithdrawPermission
		// ObjectTypeVestingBalance
		// ObjectTypeWorker
		switch id.SpaceType() {
		case types.SpaceTypeProtocol:
			switch id.ObjectType() {
			case types.ObjectTypeAccount:
				acc := types.Account{}
				if err := ffjson.Unmarshal(b, &acc); err != nil {
					return nil, errors.Annotate(err, "unmarshal Account")
				}
				ret = append(ret, acc)
			case types.ObjectTypeAsset:
				ass := types.Asset{}
				if err := ffjson.Unmarshal(b, &ass); err != nil {
					return nil, errors.Annotate(err, "unmarshal Asset")
				}
				ret = append(ret, ass)
			case types.ObjectTypeForceSettlement:
				set := types.ForceSettlementOrder{}
				if err := ffjson.Unmarshal(b, &set); err != nil {
					return nil, errors.Annotate(err, "unmarshal ForceSettlementOrder")
				}
				ret = append(ret, set)
			case types.ObjectTypeLimitOrder:
				lim := types.LimitOrder{}
				if err := ffjson.Unmarshal(b, &lim); err != nil {
					return nil, errors.Annotate(err, "unmarshal LimitOrder")
				}
				ret = append(ret, lim)
			case types.ObjectTypeCallOrder:
				cal := types.CallOrder{}
				if err := ffjson.Unmarshal(b, &cal); err != nil {
					return nil, errors.Annotate(err, "unmarshal CallOrder")
				}
				ret = append(ret, cal)
			case types.ObjectTypeCommitteeMember:
				mem := types.CommitteeMember{}
				if err := ffjson.Unmarshal(b, &mem); err != nil {
					return nil, errors.Annotate(err, "unmarshal CommitteeMember")
				}
				ret = append(ret, mem)
			case types.ObjectTypeOperationHistory:
				hist := types.OperationHistory{}
				if err := ffjson.Unmarshal(b, &hist); err != nil {
					return nil, errors.Annotate(err, "unmarshal OperationHistory")
				}
				ret = append(ret, hist)
			case types.ObjectTypeBalance:
				bal := types.Balance{}
				if err := ffjson.Unmarshal(b, &bal); err != nil {
					return nil, errors.Annotate(err, "unmarshal Balance")
				}
				ret = append(ret, bal)

			default:
				logging.DDumpUnmarshaled(id.ObjectType().String(), b)
				return nil, errors.Errorf("unable to parse GrapheneObject with ID %s", id)
			}

		case types.SpaceTypeImplementation:
			switch id.ObjectType() {
			case types.ObjectTypeAssetBitAssetData:
				bit := types.BitAssetData{}
				if err := ffjson.Unmarshal(b, &bit); err != nil {
					return nil, errors.Annotate(err, "unmarshal BitAssetData")
				}
				ret = append(ret, bit)

			default:
				logging.DDumpUnmarshaled(id.ObjectType().String(), b)
				return nil, errors.Errorf("unable to parse GrapheneObject with ID %s", id)
			}
		}
	}

	return ret, nil
}

// CancelOrder cancels an order given by orderID
func (p *websocketAPI) CancelOrder(orderID types.GrapheneObject, broadcast bool) (*types.SignedTransaction, error) {
	resp, err := p.wsClient.CallAPI(0, "cancel_order", orderID.ID(), broadcast)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("cancel_order <", resp)

	ret := types.SignedTransaction{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Transaction")
	}

	return &ret, nil
}

func (p *websocketAPI) DatabaseAPIID() int {
	return p.databaseAPIID
}

func (p *websocketAPI) BroadcastAPIID() int {
	return p.broadcastAPIID
}

func (p *websocketAPI) HistoryAPIID() int {
	return p.historyAPIID
}

//CallWsAPI invokes a websocket API call
func (p *websocketAPI) CallWsAPI(apiID int, method string, args ...interface{}) (interface{}, error) {
	return p.wsClient.CallAPI(apiID, method, args...)
}

//OnError - hook your notify callback here
func (p *websocketAPI) OnNotify(subscriberID int, notifyFn func(msg interface{}) error) error {
	return p.wsClient.OnNotify(subscriberID, notifyFn)
}

//OnError - hook your error callback here
func (p *websocketAPI) OnError(errorFn api.ErrorFunc) {
	p.wsClient.OnError(errorFn)
}

//SetCredentials defines username and password for Websocket API login.
func (p *websocketAPI) SetCredentials(username, password string) {
	p.username = username
	p.password = password
}

// Connect initializes the API and connects underlying resources
func (p *websocketAPI) Connect() error {
	if p.wsClient != nil {
		if err := p.wsClient.Connect(); err != nil {
			return errors.Annotate(err, "Connect [ws]")
		}
	}

	if ok, err := p.login(); err != nil || !ok {
		if err != nil {
			return errors.Annotate(err, "login")
		}
		return errors.New("login failed")
	}

	if err := p.getAPIIDs(); err != nil {
		return errors.Annotate(err, "getApiIDs")
	}

	chainID, err := p.GetChainID()
	if err != nil {
		return errors.Annotate(err, "GetChainID")
	}

	if err := config.SetCurrent(chainID); err != nil {
		return errors.Annotate(err, "SetCurrent")
	}

	return nil
}

func (p *websocketAPI) getAPIIDs() (err error) {
	p.databaseAPIID, err = p.getAPIID("database")
	if err != nil {
		return errors.Annotate(err, "database")
	}

	p.historyAPIID, err = p.getAPIID("history")
	if err != nil {
		return errors.Annotate(err, "history")
	}

	p.broadcastAPIID, err = p.getAPIID("network_broadcast")
	if err != nil {
		return errors.Annotate(err, "network")
	}

	return nil
}

//Close shuts the API down and closes underlying resources.
func (p *websocketAPI) Close() error {
	if p.wsClient != nil {
		if err := p.wsClient.Close(); err != nil {
			return errors.Annotate(err, "Close [ws]")
		}
		p.wsClient = nil
	}

	return nil
}

//NewWebsocketAPI creates a new WebsocketAPI interface.
//wsEndpointURL: a mandatory websocket node URL.
func NewWebsocketAPI(wsEndpointURL string) WebsocketAPI {
	api := &websocketAPI{
		databaseAPIID:  InvalidApiID,
		historyAPIID:   InvalidApiID,
		broadcastAPIID: InvalidApiID,
	}

	api.wsClient = NewSimpleClientProvider(wsEndpointURL, api)
	return api
}

//NewWebsocketAPIWithAutoEndpoint creates a new WebsocketAPI interface with automatic node latency checking.
//It's best to use this API instance type for a long API lifecycle because the latency tester takes time to unleash its magic.
//startupEndpointURL: a mandatory websocket node URL to startup the latency tester quickly.
func NewWebsocketAPIWithAutoEndpoint(startupEndpointURL string) (WebsocketAPI, error) {
	api := &websocketAPI{
		databaseAPIID:  InvalidApiID,
		historyAPIID:   InvalidApiID,
		broadcastAPIID: InvalidApiID,
	}

	pr, err := NewBestNodeClientProvider(startupEndpointURL, api)
	if err != nil {
		return nil, errors.Annotate(err, "NewBestNodeClientProvider")
	}

	api.wsClient = pr
	return api, nil
}
