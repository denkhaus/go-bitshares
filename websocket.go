package bitshares

import (
	"encoding/json"
	"time"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/config"
	"github.com/denkhaus/bitshares/crypto"
	"github.com/denkhaus/bitshares/operations"
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
	CallWsAPI(apiID int, method string, args ...interface{}) (*json.RawMessage, error)
	Close() error
	Connect() error
	DatabaseAPIID() int
	HistoryAPIID() int
	BroadcastAPIID() int
	SetCredentials(username, password string)
	OnError(api.ErrorFunc)
	Subscribe(apiID int, method string, fn api.SubscribeCallback, args ...interface{}) (*json.RawMessage, error)
	BuildSignedTransaction(keyBag *crypto.KeyBag, feeAsset types.GrapheneObject, ops ...types.Operation) (*types.SignedTransaction, error)
	SignTransaction(keyBag *crypto.KeyBag, trx *types.SignedTransaction) error

	//Websocket API functions
	BroadcastTransaction(tx *types.SignedTransaction) error
	BroadcastTransactionSynchronous(tx *types.SignedTransaction) (*types.BroadcastResponse, error)
	CancelAllSubscriptions() error
	GetAccountBalances(account types.GrapheneObject, assets ...types.GrapheneObject) (types.AssetAmounts, error)
	GetAccountByName(name string) (*types.Account, error)
	GetAccountHistory(account types.GrapheneObject, stop types.GrapheneObject, limit int, start types.GrapheneObject) (types.OperationHistories, error)
	GetAccounts(accountIDs ...types.GrapheneObject) (types.Accounts, error)
	GetBlock(number uint64) (*types.Block, error)
	GetBlockHeader(block uint64) (*types.BlockHeader, error)
	GetCallOrders(assetID types.GrapheneObject, limit int) (types.CallOrders, error)
	GetChainID() (string, error)
	GetDynamicGlobalProperties() (*types.DynamicGlobalProperties, error)
	GetForceSettlementOrders(assetID types.GrapheneObject, limit int) (types.ForceSettlementOrders, error)
	GetFullAccounts(accountIDs ...types.GrapheneObject) (types.FullAccountInfos, error)
	GetLimitOrders(base, quote types.GrapheneObject, limit int) (types.LimitOrders, error)
	GetOrderBook(base, quote types.GrapheneObject, depth int) (*types.OrderBook, error)
	GetMarginPositions(accountID types.GrapheneObject) (types.CallOrders, error)
	GetObjects(objectIDs ...types.GrapheneObject) ([]interface{}, error)
	GetPotentialSignatures(tx *types.SignedTransaction) (types.PublicKeys, error)
	GetRecentTransactionByID(transactionID uint32) (*types.SignedTransaction, error)
	GetRequiredSignatures(tx *types.SignedTransaction, keys types.PublicKeys) (types.PublicKeys, error)
	GetRequiredFees(ops types.Operations, feeAsset types.GrapheneObject) (types.AssetAmounts, error)
	GetTicker(base, quote types.GrapheneObject) (*types.MarketTicker, error)
	GetTradeHistory(base, quote types.GrapheneObject, toTime, fromTime time.Time, limit int) (types.MarketTrades, error)
	GetTransaction(blockNum uint64, trxInBlock uint32) (*types.SignedTransaction, error)
	LimitOrderCancel(keyBag *crypto.KeyBag, feePayingAccount, orderID, feeAsset types.GrapheneObject) error
	ListAssets(lowerBoundSymbol string, limit int) (types.Assets, error)
	LookupAssetSymbols(symbols ...string) (types.Assets, error)
	SetSubscribeCallback(ID uint64, clearFilter bool) error
	SubscribeToBlockApplied(onBlockApplied api.BlockAppliedCallback) error
	SubscribeToMarket(base, quote types.GrapheneObject, onMarketData api.SubscribeCallback) error
	SubscribeToPendingTransactions(onPendingTransaction api.SubscribeCallback) error
	Transfer(keyBag *crypto.KeyBag, from, to, feeAsset types.GrapheneObject, amount types.AssetAmount, memo string) error
	UnsubscribeFromMarket(base, quote types.GrapheneObject) error
	Get24Volume(base types.GrapheneObject, quote types.GrapheneObject) (*types.Volume24, error)
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
		return InvalidApiID, errors.Annotatef(err, "CallAPI %s", identifier)
	}

	logging.DDumpJSON("getApiID <", resp)

	var id int
	if err := ffjson.Unmarshal(*resp, &id); err != nil {
		return InvalidApiID, errors.Annotate(err, "Unmarshal [id]")
	}

	return id, nil
}

// login
func (p *websocketAPI) login() (bool, error) {
	resp, err := p.wsClient.CallAPI(1, "login", p.username, p.password)
	if err != nil {
		return false, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("login <", resp)

	var success bool
	if err := ffjson.Unmarshal(*resp, &success); err != nil {
		return false, errors.Annotate(err, "Unmarshal [success]")
	}

	return success, nil
}

// SetSubscribeCallback - To simplify development a global subscription callback can be registered.
// Every notification initiated by the full node will carry a particular id as defined by the user with the identifier parameter.
func (p *websocketAPI) SetSubscribeCallback(ID uint64, clearFilter bool) error {
	_, err := p.wsClient.CallAPI(p.databaseAPIID, "set_subscribe_callback", ID, clearFilter)
	if err != nil {
		return errors.Annotate(err, "CallAPI")
	}

	return nil
}

// SubscribeToPendingTransactions - Notifications for incoming unconfirmed transactions.
func (p *websocketAPI) SubscribeToPendingTransactions(onPendingTransaction api.SubscribeCallback) error {
	_, err := p.wsClient.Subscribe(p.databaseAPIID, "set_pending_transaction_callback",
		onPendingTransaction,
	)

	return err
}

// SubscribeToBlockApplied gives a notification whenever the block blockid is applied to the blockchain.
func (p *websocketAPI) SubscribeToBlockApplied(onBlockApplied api.BlockAppliedCallback) error {
	_, err := p.wsClient.Subscribe(p.databaseAPIID, "set_block_applied_callback",
		func(in interface{}) error {
			for _, id := range in.([]interface{}) {
				if err := onBlockApplied(id.(string)); err != nil {
					return err
				}
			}

			return nil
		},
	)

	return err
}

// SubscribeToMarket subscribes to market changes in market base:quote and sends notifications by callback.
func (p *websocketAPI) SubscribeToMarket(base, quote types.GrapheneObject, onMarketData api.SubscribeCallback) error {
	_, err := p.wsClient.Subscribe(p.databaseAPIID, "subscribe_to_market",
		onMarketData, base.ID(), quote.ID(),
	)

	return err
}

// UnsubscribeFromMarket
func (p *websocketAPI) UnsubscribeFromMarket(base types.GrapheneObject, quote types.GrapheneObject) error {
	// returns nil if successful
	_, err := p.wsClient.CallAPI(p.databaseAPIID, "unsubscribe_from_market", base.ID(), quote.ID())
	if err != nil {
		return errors.Annotate(err, "CallAPI")
	}

	return nil
}

// CancelAllSubscriptions
func (p *websocketAPI) CancelAllSubscriptions() error {
	// returns nil
	_, err := p.wsClient.CallAPI(p.databaseAPIID, "cancel_all_subscriptions", types.EmptyParams)
	if err != nil {
		return errors.Annotate(err, "CallAPI")
	}

	return nil
}

// BroadcastTransaction broadcasts a transaction to the network.
// The transaction will be checked for validity prior to broadcasting. If it fails to apply at the connected node,
// an error will be thrown and the transaction will not be broadcast.
func (p *websocketAPI) BroadcastTransaction(tx *types.SignedTransaction) error {
	_, err := p.wsClient.CallAPI(p.broadcastAPIID, "broadcast_transaction", tx)
	if err != nil {
		return errors.Annotate(err, "CallAPI")
	}

	return nil
}

// BroadcastTransactionSynchronous broadcasts a transaction to the network.
// The transaction will be checked for validity prior to broadcasting. If it fails to apply at the connected node,
// an error will be thrown and the transaction will not be broadcast. This version of broadcast transaction registers a callback method
// that will be called when the transaction is included into a block. The callback method includes the transaction id, block number, and transaction number in the block.
func (p *websocketAPI) BroadcastTransactionSynchronous(tx *types.SignedTransaction) (*types.BroadcastResponse, error) {
	resp, err := p.wsClient.CallAPI(p.broadcastAPIID, "broadcast_transaction_synchronous", tx)
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	var ret types.BroadcastResponse
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [BroadcastResponse]")
	}

	return &ret, nil
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
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_potential_signatures <", resp)

	ret := types.PublicKeys{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [PublicKeys]")
	}

	return ret, nil
}

// GetTransaction used to fetch an individual transaction.
func (p *websocketAPI) GetTransaction(blockNum uint64, trxInBlock uint32) (*types.SignedTransaction, error) {
	resp, err := p.wsClient.CallAPI(p.databaseAPIID, "get_transaction", blockNum, trxInBlock)
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_transaction <", resp)

	ret := types.SignedTransaction{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [Transaction]")
	}

	return &ret, nil
}

// GetRecentTransactionByID
// If the transaction has not expired, this method will return the transaction for the given ID or
// it will return nil if it is not known. Just because it is not known does not mean
// it wasn’t included in the blockchain.
func (p *websocketAPI) GetRecentTransactionByID(transactionID uint32) (*types.SignedTransaction, error) {
	resp, err := p.wsClient.CallAPI(p.databaseAPIID, "get_recent_transaction_by_id", transactionID)
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_recent_transaction_by_id <", resp)

	ret := types.SignedTransaction{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [Transaction]")
	}

	return &ret, nil
}

//GetRequiredSignatures returns the minimum subset of public keys to sign a transaction.
func (p *websocketAPI) GetRequiredSignatures(tx *types.SignedTransaction, potKeys types.PublicKeys) (types.PublicKeys, error) {
	resp, err := p.wsClient.CallAPI(p.databaseAPIID, "get_required_signatures", tx, potKeys)
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_required_signatures <", resp)

	ret := types.PublicKeys{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [PublicKeys]")
	}

	return ret, nil
}

//GetBlock returns a Block by number.
func (p *websocketAPI) GetBlock(block uint64) (*types.Block, error) {
	resp, err := p.wsClient.CallAPI(0, "get_block", block)
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_block <", resp)

	ret := types.Block{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [Block]")
	}

	return &ret, nil
}

// GetBlockHeader returns block header by number.
func (p *websocketAPI) GetBlockHeader(block uint64) (*types.BlockHeader, error) {
	resp, err := p.wsClient.CallAPI(0, "get_block_header", block)
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_block_header <", resp)

	ret := types.BlockHeader{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [BlockHeader]")
	}

	return &ret, nil
}

// GetTicker returns the ticker for the market base:quote for the last 24 h
func (p *websocketAPI) GetTicker(base, quote types.GrapheneObject) (*types.MarketTicker, error) {
	resp, err := p.wsClient.CallAPI(0, "get_ticker", base.ID(), quote.ID())
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_ticker <", resp)

	ret := types.MarketTicker{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [MarketTicker]")
	}

	return &ret, nil
}

//GetAccountByName returns a Account object by username
func (p *websocketAPI) GetAccountByName(name string) (*types.Account, error) {
	resp, err := p.wsClient.CallAPI(0, "get_account_by_name", name)
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_account_by_name <", resp)

	ret := types.Account{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [Account]")
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
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_account_history <", resp)

	ret := types.OperationHistories{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [Histories]")
	}

	return ret, nil
}

//GetAccounts returns a list of accounts by accountID(s).
func (p *websocketAPI) GetAccounts(accounts ...types.GrapheneObject) (types.Accounts, error) {
	ids := types.GrapheneObjects(accounts).ToStrings()
	resp, err := p.wsClient.CallAPI(0, "get_accounts", ids)
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_accounts <", resp)

	ret := types.Accounts{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [Accounts]")
	}

	return ret, nil
}

//GetDynamicGlobalProperties returns essential runtime properties of bitshares network
func (p *websocketAPI) GetDynamicGlobalProperties() (*types.DynamicGlobalProperties, error) {
	resp, err := p.wsClient.CallAPI(0, "get_dynamic_global_properties", types.EmptyParams)
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_dynamic_global_properties <", resp)

	ret := types.DynamicGlobalProperties{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [DynamicGlobalProperties]")
	}

	return &ret, nil
}

//GetAccountBalances retrieves AssetAmounts by given AccountID
func (p *websocketAPI) GetAccountBalances(account types.GrapheneObject, assets ...types.GrapheneObject) (types.AssetAmounts, error) {
	ids := types.GrapheneObjects(assets).ToStrings()
	resp, err := p.wsClient.CallAPI(0, "get_account_balances", account.ID(), ids)
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_account_balances <", resp)

	ret := types.AssetAmounts{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [AssetAmounts]")
	}

	return ret, nil
}

// GetFullAccounts retrieves full account information by given AccountIDs
func (p *websocketAPI) GetFullAccounts(accounts ...types.GrapheneObject) (types.FullAccountInfos, error) {
	ids := types.GrapheneObjects(accounts).ToStrings()
	resp, err := p.wsClient.CallAPI(0, "get_full_accounts", ids, false) //do not subscribe for now
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_full_accounts <", resp)

	ret := types.FullAccountInfos{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [FullAccountInfos]")
	}

	return ret, nil
}

// Get24Volume returns the base:quote assets 24h volume
func (p *websocketAPI) Get24Volume(base, quote types.GrapheneObject) (*types.Volume24, error) {
	resp, err := p.wsClient.CallAPI(p.databaseAPIID, "get_24_volume", base.ID(), quote.ID())
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_24_volume <", resp)

	ret := types.Volume24{}
	if err = ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [Volume24]")
	}

	return &ret, nil
}

// ListAssets retrieves assets
// lowerBoundSymbol: Lower bound of symbol names to retrieve
// limit: Maximum number of assets to fetch, if the constant AssetsListAll is passed, all existing assets will be retrieved.
func (p *websocketAPI) ListAssets(lowerBoundSymbol string, limit int) (types.Assets, error) {
	if limit > AssetsMaxBatchSize {
		limit = AssetsMaxBatchSize
	}

	resp, err := p.wsClient.CallAPI(0, "list_assets", lowerBoundSymbol, limit)
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("list_assets <", resp)

	ret := types.Assets{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [Assets]")
	}

	return ret, nil
}

// LookupAssetSymbols get assets corresponding to the provided symbols or IDs
func (p *websocketAPI) LookupAssetSymbols(symbols ...string) (types.Assets, error) {
	resp, err := p.wsClient.CallAPI(0, "lookup_asset_symbols", symbols)
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("lookup_asset_symbols <", resp)

	ret := types.Assets{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [Assets]")
	}

	return ret, nil
}

//GetRequiredFees calculates the required fee for each operation by the specified asset type.
func (p *websocketAPI) GetRequiredFees(ops types.Operations, feeAsset types.GrapheneObject) (types.AssetAmounts, error) {
	resp, err := p.wsClient.CallAPI(0, "get_required_fees", ops.Envelopes(), feeAsset.ID())
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_required_fees <", resp)

	ret := types.AssetAmounts{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [AssetAmounts]")
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
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_limit_orders <", resp)

	ret := types.LimitOrders{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [LimitOrders]")
	}

	return ret, nil
}

// LimitOrderCancel cancels a certain limit order given by orderID. Fees are paid in feeAsset.
// The transaction is signed with private keys in keyBag.
func (p *websocketAPI) LimitOrderCancel(keyBag *crypto.KeyBag, feePayingAccount, orderID, feeAsset types.GrapheneObject) error {
	op := operations.LimitOrderCancelOperation{
		FeePayingAccount: types.AccountIDFromObject(feePayingAccount),
		Order:            types.LimitOrderIDFromObject(orderID),
		Extensions:       types.Extensions{},
	}

	trx, err := p.BuildSignedTransaction(keyBag, feeAsset, &op)
	if err != nil {
		return errors.Annotate(err, "BuildSignedTransaction")
	}

	if err := p.BroadcastTransaction(trx); err != nil {
		return errors.Annotate(err, "BroadcastTransaction")
	}

	return nil
}

//GetOrderBook returns the OrderBook for the market base:quote.
func (p *websocketAPI) GetOrderBook(base, quote types.GrapheneObject, depth int) (*types.OrderBook, error) {
	resp, err := p.wsClient.CallAPI(0, "get_order_book", base.ID(), quote.ID(), depth)
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_order_book <", resp)

	ret := types.OrderBook{}
	if err = ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [OrderBook]")
	}

	return &ret, nil
}

//GetForceSettlementOrders returns ForceSettlementOrders type.
func (p *websocketAPI) GetForceSettlementOrders(assetID types.GrapheneObject, limit int) (types.ForceSettlementOrders, error) {
	if limit > GetForceSettlementOrdersLimit {
		limit = GetForceSettlementOrdersLimit
	}

	resp, err := p.wsClient.CallAPI(0, "get_settle_orders", assetID.ID(), limit)
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_settle_orders <", resp)

	ret := types.ForceSettlementOrders{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [ForceSettlementOrders]")
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
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_call_orders <", resp)

	ret := types.CallOrders{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [CallOrders]")
	}

	return ret, nil
}

//GetMarginPositions returns CallOrders type.
func (p *websocketAPI) GetMarginPositions(accountID types.GrapheneObject) (types.CallOrders, error) {
	resp, err := p.wsClient.CallAPI(0, "get_margin_positions", accountID.ID())
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_margin_positions <", resp)

	ret := types.CallOrders{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [CallOrders]")
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
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_trade_history <", resp)

	ret := types.MarketTrades{}
	if err := ffjson.Unmarshal(*resp, &ret); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [MarketTrades]")
	}

	return ret, nil
}

//GetChainID returns the ID of the chain we are connected to.
func (p *websocketAPI) GetChainID() (string, error) {
	resp, err := p.wsClient.CallAPI(p.databaseAPIID, "get_chain_id", types.EmptyParams)
	if err != nil {
		return "", errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_chain_id <", resp)

	var id string
	if err := ffjson.Unmarshal(*resp, &id); err != nil {
		return "", errors.Annotate(err, "Unmarshal [id]")
	}

	return id, nil
}

//GetObjects returns a list of Graphene Objects by ID.
func (p *websocketAPI) GetObjects(ids ...types.GrapheneObject) ([]interface{}, error) {
	params := types.GrapheneObjects(ids).ToStrings()
	resp, err := p.wsClient.CallAPI(0, "get_objects", params)
	if err != nil {
		return nil, errors.Annotate(err, "CallAPI")
	}

	logging.DDumpJSON("get_objects <", resp)

	var data []interface{}
	if err := ffjson.Unmarshal(*resp, &data); err != nil {
		return nil, errors.Annotate(err, "Unmarshal [data]")
	}

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
		// ObjectTypeWorker
		switch id.SpaceType() {
		case types.SpaceTypeProtocol:
			switch id.ObjectType() {
			case types.ObjectTypeVestingBalance:
				t := types.VestingBalance{}
				if err := t.UnmarshalJSON(b); err != nil {
					return nil, errors.Annotate(err, "Unmarshal [VestingBalance]")
				}
				ret = append(ret, t)
			case types.ObjectTypeAccount:
				t := types.Account{}
				if err := t.UnmarshalJSON(b); err != nil {
					return nil, errors.Annotate(err, "Unmarshal [Account]")
				}
				ret = append(ret, t)
			case types.ObjectTypeAsset:
				t := types.Asset{}
				if err := t.UnmarshalJSON(b); err != nil {
					return nil, errors.Annotate(err, "Unmarshal [Asset]")
				}
				ret = append(ret, t)
			case types.ObjectTypeForceSettlement:
				t := types.ForceSettlementOrder{}
				if err := t.UnmarshalJSON(b); err != nil {
					return nil, errors.Annotate(err, "Unmarshal [ForceSettlementOrder]")
				}
				ret = append(ret, t)
			case types.ObjectTypeLimitOrder:
				t := types.LimitOrder{}
				if err := t.UnmarshalJSON(b); err != nil {
					return nil, errors.Annotate(err, "Unmarshal [LimitOrder]")
				}
				ret = append(ret, t)
			case types.ObjectTypeCallOrder:
				t := types.CallOrder{}
				if err := t.UnmarshalJSON(b); err != nil {
					return nil, errors.Annotate(err, "Unmarshal [CallOrder]")
				}
				ret = append(ret, t)
			case types.ObjectTypeCommitteeMember:
				t := types.CommitteeMember{}
				if err := t.UnmarshalJSON(b); err != nil {
					return nil, errors.Annotate(err, "Unmarshal [CommitteeMember]")
				}
				ret = append(ret, t)
			case types.ObjectTypeOperationHistory:
				t := types.OperationHistory{}
				if err := t.UnmarshalJSON(b); err != nil {
					return nil, errors.Annotate(err, "Unmarshal [OperationHistory]")
				}
				ret = append(ret, t)
			case types.ObjectTypeBalance:
				t := types.Balance{}
				if err := t.UnmarshalJSON(b); err != nil {
					return nil, errors.Annotate(err, "Unmarshal [Balance]")
				}
				ret = append(ret, t)

			default:
				logging.DDumpUnmarshaled(id.ObjectType().String(), b)
				return nil, errors.Errorf("unable to parse Object with ID %s", id)
			}

			// TODO: implement
			// ObjectTypeGlobalProperty
			// ObjectTypeAssetDynamicData
			// ObjectTypeBlockSummary
			// ObjectTypeAccountTransactionHistory
			// ObjectTypeBlindedBalance
			// ObjectTypeChainProperty
			// ObjectTypeWitnessSchedule
			// ObjectTypeBudgetRecord
		case types.SpaceTypeImplementation:
			switch id.ObjectType() {
			case types.ObjectTypeSpecialAuthority:
				t := types.SpecialAuthority{}
				if err := t.UnmarshalJSON(b); err != nil {
					return nil, errors.Annotate(err, "Unmarshal [SpecialAuthority]")
				}
				ret = append(ret, t)
			case types.ObjectTypeTransaction:
				t := types.Transaction{}
				if err := t.UnmarshalJSON(b); err != nil {
					return nil, errors.Annotate(err, "Unmarshal [Transaction]")
				}
				ret = append(ret, t)
			case types.ObjectTypeDynamicGlobalProperty:
				t := types.DynamicGlobalProperties{}
				if err := t.UnmarshalJSON(b); err != nil {
					return nil, errors.Annotate(err, "Unmarshal [DynamicGlobalProperties]")
				}
				ret = append(ret, t)
			case types.ObjectTypeAccountStatistics:
				t := types.AccountStatistics{}
				if err := t.UnmarshalJSON(b); err != nil {
					return nil, errors.Annotate(err, "Unmarshal [AccountStatistics]")
				}
				ret = append(ret, t)
			case types.ObjectTypeAccountBalance:
				t := types.AccountBalance{}
				if err := t.UnmarshalJSON(b); err != nil {
					return nil, errors.Annotate(err, "Unmarshal [AccountBalance]")
				}
				ret = append(ret, t)
			case types.ObjectTypeAssetBitAssetData:
				t := types.BitAssetData{}
				if err := t.UnmarshalJSON(b); err != nil {
					return nil, errors.Annotate(err, "Unmarshal [BitAssetData]")
				}
				ret = append(ret, t)

			default:
				logging.DDumpUnmarshaled(id.ObjectType().String(), b)
				return nil, errors.Errorf("unable to parse Object with ID %s", id)
			}
		}
	}

	return ret, nil
}

// Transfer transfers a certain amount between two accounts. Fees are paid in feeAsset.
// The transaction is signed with private keys in keyBag.
func (p *websocketAPI) Transfer(keyBag *crypto.KeyBag, from, to, feeAsset types.GrapheneObject, amount types.AssetAmount, memo string) error {
	op := operations.TransferOperation{
		Amount:     amount,
		Extensions: types.Extensions{},
		From:       types.AccountIDFromObject(from),
		To:         types.AccountIDFromObject(to),
	}

	if memo != "" {
		builder := p.NewMemoBuilder(from, to, memo)
		m, err := builder.Encrypt(keyBag)
		if err != nil {
			return errors.Annotate(err, "Encrypt [memo]")
		}

		op.Memo = m
	}

	trx, err := p.BuildSignedTransaction(keyBag, feeAsset, &op)
	if err != nil {
		return errors.Annotate(err, "BuildSignedTransaction")
	}

	if err := p.BroadcastTransaction(trx); err != nil {
		return errors.Annotate(err, "BroadcastTransaction")
	}

	return nil
}

//DatabaseAPIID returns the database API ID
func (p *websocketAPI) DatabaseAPIID() int {
	return p.databaseAPIID
}

//BroadcastAPIID returns the broadcast API ID
func (p *websocketAPI) BroadcastAPIID() int {
	return p.broadcastAPIID
}

//HistoryAPIID returns the history API ID
func (p *websocketAPI) HistoryAPIID() int {
	return p.historyAPIID
}

//CallWsAPI invokes a websocket API call
func (p *websocketAPI) CallWsAPI(apiID int, method string, args ...interface{}) (*json.RawMessage, error) {
	return p.wsClient.CallAPI(apiID, method, args...)
}

//Subscribe - hook your subscribe callback here
func (p *websocketAPI) Subscribe(apiID int, method string, fn api.SubscribeCallback, args ...interface{}) (*json.RawMessage, error) {
	return p.wsClient.Subscribe(apiID, method, fn, args...)
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
//wsEndpointURL: a websocket node endpoint URL.
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
//startupEndpointURL: a websocket node endpoint URL to startup the latency tester quickly.
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
