package api

import (
	"time"

	"github.com/denkhaus/bitshares/client"
	"github.com/denkhaus/bitshares/config"
	"github.com/denkhaus/bitshares/crypto"
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

const (
	InvalidApiID           = -1
	AssetsListAll          = -1
	AssetsMaxBatchSize     = 100
	GetCallOrdersLimit     = 100
	GetLimitOrdersLimit    = 100
	GetSettleOrdersLimit   = 100
	GetTradeHistoryLimit   = 100
	GetAccountHistoryLimit = 100
)

type BitsharesAPI interface {

	//Common functions
	SetDebug(debug bool)
	Debug(descr string, in interface{})
	CallWsAPI(apiID int, method string, args ...interface{}) (interface{}, error)
	Close() error
	Connect() error
	DatabaseAPIID() int
	CryptoAPIID() int
	HistoryAPIID() int
	BroadcastAPIID() int
	SetCredentials(username, password string)
	OnError(client.ErrorFunc)
	OnNotify(subscriberID int, notifyFn func(msg interface{}) error) error
	BuildSignedTransaction(keyBag *crypto.KeyBag, feeAsset types.GrapheneObject, ops ...types.Operation) (*types.SignedTransaction, error)
	VerifySignedTransaction(keyBag *crypto.KeyBag, tx *types.SignedTransaction) (bool, error)
	SignTransaction(keyBag *crypto.KeyBag, trx *types.SignedTransaction) error
	SignWithKeys(types.PrivateKeys, *types.SignedTransaction) error

	//Websocket API functions
	BroadcastTransaction(tx *types.SignedTransaction) error
	CancelAllSubscriptions() error
	CancelOrder(orderID types.GrapheneObject, broadcast bool) (*types.SignedTransaction, error)
	GetAccountBalances(account types.GrapheneObject, assets ...types.GrapheneObject) (types.AssetAmounts, error)
	GetAccountByName(name string) (*types.Account, error)
	GetAccountHistory(account types.GrapheneObject, stop types.GrapheneObject, limit int, start types.GrapheneObject) (types.OperationHistories, error)
	GetAccounts(accountIDs ...types.GrapheneObject) (types.Accounts, error)
	GetBlock(number uint64) (*types.Block, error)
	GetCallOrders(assetID types.GrapheneObject, limit int) (types.CallOrders, error)
	GetChainID() (string, error)
	GetDynamicGlobalProperties() (*types.DynamicGlobalProperties, error)
	GetLimitOrders(base, quote types.GrapheneObject, limit int) (types.LimitOrders, error)
	GetMarginPositions(accountID types.GrapheneObject) (types.CallOrders, error)
	GetObjects(objectIDs ...types.GrapheneObject) ([]interface{}, error)
	GetPotentialSignatures(tx *types.SignedTransaction) (types.PublicKeys, error)
	GetRequiredSignatures(tx *types.SignedTransaction, keys types.PublicKeys) (types.PublicKeys, error)
	GetRequiredFees(ops types.Operations, feeAsset types.GrapheneObject) (types.AssetAmounts, error)
	GetSettleOrders(assetID types.GrapheneObject, limit int) (types.SettleOrders, error)
	GetTradeHistory(base, quote types.GrapheneObject, toTime, fromTime time.Time, limit int) (types.MarketTrades, error)
	ListAssets(lowerBoundSymbol string, limit int) (types.Assets, error)
	SetSubscribeCallback(notifyID int, clearFilter bool) error
	SubscribeToMarket(notifyID int, base types.GrapheneObject, quote types.GrapheneObject) error
	UnsubscribeFromMarket(base types.GrapheneObject, quote types.GrapheneObject) error

	//Wallet API functions
	WalletListAccountBalances(account types.GrapheneObject) (types.AssetAmounts, error)
	WalletLock() error
	WalletUnlock(password string) error
	WalletIsLocked() (bool, error)
	WalletBorrowAsset(account types.GrapheneObject, amountToBorrow string, symbolToBorrow types.GrapheneObject, amountOfCollateral string, broadcast bool) (*types.SignedTransaction, error)
	WalletBuy(account types.GrapheneObject, base, quote types.GrapheneObject, rate string, amount string, broadcast bool) (*types.SignedTransaction, error)
	WalletBuyEx(account types.GrapheneObject, base, quote types.GrapheneObject, rate float64, amount float64, broadcast bool) (*types.SignedTransaction, error)
	WalletSell(account types.GrapheneObject, base, quote types.GrapheneObject, rate string, amount string, broadcast bool) (*types.SignedTransaction, error)
	WalletSellEx(account types.GrapheneObject, base, quote types.GrapheneObject, rate float64, amount float64, broadcast bool) (*types.SignedTransaction, error)
	WalletSellAsset(account types.GrapheneObject, amountToSell string, symbolToSell types.GrapheneObject, minToReceive string, symbolToReceive types.GrapheneObject, timeout uint32, fillOrKill bool, broadcast bool) (*types.SignedTransaction, error)
	WalletSignTransaction(tx *types.SignedTransaction, broadcast bool) (*types.SignedTransaction, error)
	WalletSerializeTransaction(tx *types.SignedTransaction) (string, error)
}

type bitsharesAPI struct {
	wsClient       ClientProvider
	rpcClient      client.RPCClient
	username       string
	password       string
	databaseAPIID  int
	historyAPIID   int
	cryptoAPIID    int
	broadcastAPIID int
	debug          bool
}

func (p *bitsharesAPI) getAPIID(identifier string) (int, error) {
	defer p.SetDebug(false)
	resp, err := p.wsClient.CallAPI(1, identifier, types.EmptyParams)
	if err != nil {
		return InvalidApiID, errors.Annotate(err, identifier)
	}

	p.Debug("getApiID <", resp)

	return int(resp.(float64)), nil
}

// login
func (p *bitsharesAPI) login() (bool, error) {
	defer p.SetDebug(false)
	resp, err := p.wsClient.CallAPI(1, "login", p.username, p.password)
	if err != nil {
		return false, err
	}

	p.Debug("login <", resp)

	return resp.(bool), nil
}

// SetSubscribeCallback
func (p *bitsharesAPI) SetSubscribeCallback(notifyID int, clearFilter bool) error {
	defer p.SetDebug(false)
	// returns nil if successfull
	_, err := p.wsClient.CallAPI(p.databaseAPIID, "set_subscribe_callback", notifyID, clearFilter)
	if err != nil {
		return err
	}

	return nil
}

// SubscribeToMarket
func (p *bitsharesAPI) SubscribeToMarket(notifyID int, base types.GrapheneObject, quote types.GrapheneObject) error {
	defer p.SetDebug(false)
	// returns nil if successfull
	_, err := p.wsClient.CallAPI(p.databaseAPIID, "subscribe_to_market", notifyID, base.ID(), quote.ID())
	if err != nil {
		return err
	}

	return nil
}

// UnsubscribeFromMarket
func (p *bitsharesAPI) UnsubscribeFromMarket(base types.GrapheneObject, quote types.GrapheneObject) error {
	defer p.SetDebug(false)
	// returns nil if successfull
	_, err := p.wsClient.CallAPI(p.databaseAPIID, "unsubscribe_from_market", base.ID(), quote.ID())
	if err != nil {
		return err
	}

	return nil
}

// CancelAllSubscriptions
func (p *bitsharesAPI) CancelAllSubscriptions() error {
	defer p.SetDebug(false)
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
func (p *bitsharesAPI) BroadcastTransaction(tx *types.SignedTransaction) error {
	defer p.SetDebug(false)

	_, err := p.wsClient.CallAPI(p.broadcastAPIID, "broadcast_transaction", tx)
	if err != nil {
		return err
	}

	return nil
}

//SignTransaction signs a given transaction.
//Required signing keys get selected by API and have to be in keyBag.
func (p *bitsharesAPI) SignTransaction(keyBag *crypto.KeyBag, tx *types.SignedTransaction) error {
	defer p.SetDebug(false)

	reqPk, err := p.RequiredSigningKeys(tx)
	if err != nil {
		return errors.Annotate(err, "RequiredSigningKeys")
	}

	signer := crypto.NewTransactionSigner(tx)

	privKeys := keyBag.PrivatesByPublics(reqPk)
	if len(privKeys) == 0 {
		return types.ErrNoSigningKeyFound
	}

	if err := signer.Sign(privKeys, config.CurrentConfig()); err != nil {
		return errors.Annotate(err, "Sign")
	}

	return nil
}

//SignWithKeys signs a given transaction with given private keys.
func (p *bitsharesAPI) SignWithKeys(keys types.PrivateKeys, tx *types.SignedTransaction) error {
	defer p.SetDebug(false)

	signer := crypto.NewTransactionSigner(tx)
	if err := signer.Sign(keys, config.CurrentConfig()); err != nil {
		return errors.Annotate(err, "Sign")
	}

	return nil
}

//VerifySignedTransaction verifies a signed transaction against all available keys in keyBag.
//If all required key are found the function returns true, otherwise false.
func (p *bitsharesAPI) VerifySignedTransaction(keyBag *crypto.KeyBag, tx *types.SignedTransaction) (bool, error) {
	defer p.SetDebug(false)

	signer := crypto.NewTransactionSigner(tx)
	verified, err := signer.Verify(keyBag, config.CurrentConfig())
	if err != nil {
		return false, errors.Annotate(err, "Verify")
	}

	return verified, nil
}

//BuildSignedTransaction builds a new transaction by given operation(s),
//applies fees, current block data and signs the transaction.
func (p *bitsharesAPI) BuildSignedTransaction(keyBag *crypto.KeyBag, feeAsset types.GrapheneObject, ops ...types.Operation) (*types.SignedTransaction, error) {
	defer p.SetDebug(false)

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

	if err := signer.Sign(privKeys, config.CurrentConfig()); err != nil {
		return nil, errors.Annotate(err, "Sign")
	}

	return tx, nil
}

//RequiredSigningKeys is a convenience call to retrieve the minimum subset of public keys to sign a transaction.
//If the transaction is already signed, the result is empty.
func (p *bitsharesAPI) RequiredSigningKeys(tx *types.SignedTransaction) (types.PublicKeys, error) {
	defer p.SetDebug(false)

	potPk, err := p.GetPotentialSignatures(tx)
	if err != nil {
		return nil, errors.Annotate(err, "GetPotentialSignatures")
	}

	p.Debug("potential pubkeys <", potPk)

	reqPk, err := p.GetRequiredSignatures(tx, potPk)
	if err != nil {
		return nil, errors.Annotate(err, "GetRequiredSignatures")
	}

	p.Debug("required pubkeys <", reqPk)

	return reqPk, nil
}

//GetPotentialSignatures will return the set of all public keys that could possibly sign for a given transaction.
//This call can be used by wallets to filter their set of public keys to just the relevant subset prior to calling
//GetRequiredSignatures to get the minimum subset.
func (p *bitsharesAPI) GetPotentialSignatures(tx *types.SignedTransaction) (types.PublicKeys, error) {
	defer p.SetDebug(false)

	resp, err := p.wsClient.CallAPI(p.databaseAPIID, "get_potential_signatures", tx)
	if err != nil {
		return nil, err
	}

	p.Debug("get_potential_signatures <", resp)

	ret := types.PublicKeys{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal PublicKeys")
	}

	return ret, nil
}

//GetRequiredSignatures returns the minimum subset of public keys to sign a transaction.
func (p *bitsharesAPI) GetRequiredSignatures(tx *types.SignedTransaction, potKeys types.PublicKeys) (types.PublicKeys, error) {
	defer p.SetDebug(false)

	resp, err := p.wsClient.CallAPI(p.databaseAPIID, "get_required_signatures", tx, potKeys)
	if err != nil {
		return nil, err
	}

	p.Debug("get_required_signatures <", resp)

	ret := types.PublicKeys{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal PublicKeys")
	}

	return ret, nil
}

//GetBlock returns a Block by block number.
func (p *bitsharesAPI) GetBlock(number uint64) (*types.Block, error) {
	defer p.SetDebug(false)
	resp, err := p.wsClient.CallAPI(0, "get_block", number)
	if err != nil {
		return nil, err
	}

	p.Debug("get_block <", resp)

	ret := types.Block{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Block")
	}

	return &ret, nil
}

//GetAccountByName returns a Account object by username
func (p *bitsharesAPI) GetAccountByName(name string) (*types.Account, error) {
	defer p.SetDebug(false)
	resp, err := p.wsClient.CallAPI(0, "get_account_by_name", name)
	if err != nil {
		return nil, err
	}

	p.Debug("get_account_by_name <", resp)

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
func (p *bitsharesAPI) GetAccountHistory(account types.GrapheneObject, stop types.GrapheneObject, limit int, start types.GrapheneObject) (types.OperationHistories, error) {
	defer p.SetDebug(false)

	if limit > GetAccountHistoryLimit {
		limit = GetAccountHistoryLimit
	}

	resp, err := p.wsClient.CallAPI(p.historyAPIID, "get_account_history", account.ID(), stop.ID(), limit, start.ID())
	if err != nil {
		return nil, err
	}

	p.Debug("get_account_history <", resp)

	ret := types.OperationHistories{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Histories")
	}

	return ret, nil
}

//GetAccounts returns a list of accounts by accountID(s).
func (p *bitsharesAPI) GetAccounts(accounts ...types.GrapheneObject) (types.Accounts, error) {
	defer p.SetDebug(false)

	ids := types.GrapheneObjects(accounts).ToStrings()
	resp, err := p.wsClient.CallAPI(0, "get_accounts", ids)
	if err != nil {
		return nil, err
	}

	p.Debug("get_accounts <", resp)

	ret := types.Accounts{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Accounts")
	}

	return ret, nil
}

//GetDynamicGlobalProperties returns essential runtime properties of bitshares network.
func (p *bitsharesAPI) GetDynamicGlobalProperties() (*types.DynamicGlobalProperties, error) {
	defer p.SetDebug(false)

	resp, err := p.wsClient.CallAPI(0, "get_dynamic_global_properties", types.EmptyParams)
	if err != nil {
		return nil, err
	}

	p.Debug("get_dynamic_global_properties <", resp)

	var ret types.DynamicGlobalProperties
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal DynamicGlobalProperties")
	}

	return &ret, nil
}

//GetAccountBalances retrieves AssetAmounts by given AccountID
func (p *bitsharesAPI) GetAccountBalances(account types.GrapheneObject, assets ...types.GrapheneObject) (types.AssetAmounts, error) {
	defer p.SetDebug(false)

	ids := types.GrapheneObjects(assets).ToStrings()
	resp, err := p.wsClient.CallAPI(0, "get_account_balances", account.ID(), ids)
	if err != nil {
		return nil, err
	}

	p.Debug("get_account_balances <", resp)

	ret := types.AssetAmounts{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal AssetAmounts")
	}

	return ret, nil
}

//ListAssets retrieves assets
//lowerBoundSymbol: Lower bound of symbol names to retrieve
//limit: Maximum number of assets to fetch, if the constant AssetsListAll is passed, all existing assets will be retrieved.
func (p *bitsharesAPI) ListAssets(lowerBoundSymbol string, limit int) (types.Assets, error) {
	defer p.SetDebug(false)

	if limit > AssetsMaxBatchSize || limit == AssetsListAll {
		limit = AssetsMaxBatchSize
	}

	resp, err := p.wsClient.CallAPI(0, "list_assets", lowerBoundSymbol, limit)
	if err != nil {
		return nil, err
	}

	p.Debug("list_assets <", resp)

	ret := types.Assets{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Assets")
	}

	return ret, nil
}

//GetRequiredFees calculates the required fee for each operation by the specified asset type.
func (p *bitsharesAPI) GetRequiredFees(ops types.Operations, feeAsset types.GrapheneObject) (types.AssetAmounts, error) {
	defer p.SetDebug(false)

	resp, err := p.wsClient.CallAPI(0, "get_required_fees", ops.Types(), feeAsset.ID())
	if err != nil {
		return nil, err
	}

	p.Debug("get_required_fees <", resp)

	ret := types.AssetAmounts{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal AssetAmounts")
	}

	return ret, nil
}

//GetLimitOrders returns LimitOrders type.
func (p *bitsharesAPI) GetLimitOrders(base, quote types.GrapheneObject, limit int) (types.LimitOrders, error) {
	defer p.SetDebug(false)

	if limit > GetLimitOrdersLimit {
		limit = GetLimitOrdersLimit
	}

	resp, err := p.wsClient.CallAPI(0, "get_limit_orders", base.ID(), quote.ID(), limit)
	if err != nil {
		return nil, err
	}

	p.Debug("get_limit_orders <", resp)

	ret := types.LimitOrders{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal LimitOrders")
	}

	return ret, nil
}

//GetSettleOrders returns SettleOrders type.
func (p *bitsharesAPI) GetSettleOrders(assetID types.GrapheneObject, limit int) (types.SettleOrders, error) {
	defer p.SetDebug(false)

	if limit > GetSettleOrdersLimit {
		limit = GetSettleOrdersLimit
	}

	resp, err := p.wsClient.CallAPI(0, "get_settle_orders", assetID.ID(), limit)
	if err != nil {
		return nil, err
	}

	p.Debug("get_settle_orders <", resp)

	ret := types.SettleOrders{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal SettleOrders")
	}

	return ret, nil
}

//GetCallOrders returns CallOrders type.
func (p *bitsharesAPI) GetCallOrders(assetID types.GrapheneObject, limit int) (types.CallOrders, error) {
	defer p.SetDebug(false)

	if limit > GetCallOrdersLimit {
		limit = GetCallOrdersLimit
	}

	resp, err := p.wsClient.CallAPI(0, "get_call_orders", assetID.ID(), limit)
	if err != nil {
		return nil, err
	}

	p.Debug("get_call_orders <", resp)

	ret := types.CallOrders{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal CallOrders")
	}

	return ret, nil
}

//GetMarginPositions returns CallOrders type.
func (p *bitsharesAPI) GetMarginPositions(accountID types.GrapheneObject) (types.CallOrders, error) {
	defer p.SetDebug(false)

	resp, err := p.wsClient.CallAPI(0, "get_margin_positions", accountID.ID())
	if err != nil {
		return nil, err
	}

	p.Debug("get_margin_positions <", resp)

	ret := types.CallOrders{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal CallOrders")
	}

	return ret, nil
}

//GetTradeHistory returns MarketTrades type.
func (p *bitsharesAPI) GetTradeHistory(base, quote types.GrapheneObject, toTime, fromTime time.Time, limit int) (types.MarketTrades, error) {
	defer p.SetDebug(false)

	if limit > GetTradeHistoryLimit {
		limit = GetTradeHistoryLimit
	}

	resp, err := p.wsClient.CallAPI(0, "get_trade_history", base.ID(), quote.ID(), toTime, fromTime, limit)
	if err != nil {
		return nil, err
	}

	p.Debug("get_trade_history <", resp)

	ret := types.MarketTrades{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal MarketTrades")
	}

	return ret, nil
}

//GetChainID returns the ID of the chain we are connected to.
func (p *bitsharesAPI) GetChainID() (string, error) {
	defer p.SetDebug(false)

	resp, err := p.wsClient.CallAPI(p.databaseAPIID, "get_chain_id", types.EmptyParams)
	if err != nil {
		return "", err
	}

	p.Debug("get_chain_id <", resp)

	return resp.(string), nil
}

//GetObjects returns a list of Graphene Objects by ID.
func (p *bitsharesAPI) GetObjects(ids ...types.GrapheneObject) ([]interface{}, error) {
	defer p.SetDebug(false)

	params := types.GrapheneObjects(ids).ToStrings()
	resp, err := p.wsClient.CallAPI(0, "get_objects", params)
	if err != nil {
		return nil, err
	}

	p.Debug("get_objects <", resp)

	data := resp.([]interface{})
	ret := make([]interface{}, 0)
	id := types.GrapheneID{}

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
		switch id.Space() {
		case types.SpaceTypeProtocol:
			switch id.Type() {
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
				set := types.SettleOrder{}
				if err := ffjson.Unmarshal(b, &set); err != nil {
					return nil, errors.Annotate(err, "unmarshal SettleOrder")
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
			case types.ObjectTypeCommiteeMember:
				mem := types.CommiteeMember{}
				if err := ffjson.Unmarshal(b, &mem); err != nil {
					return nil, errors.Annotate(err, "unmarshal CommiteeMember")
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
				util.DumpUnmarshaled(id.Type().String(), b)
				return nil, errors.Errorf("unable to parse GrapheneObject with ID %s", id)
			}

		case types.SpaceTypeImplementation:
			switch id.Type() {
			case types.ObjectTypeAssetBitAssetData:
				bit := types.BitAssetData{}
				if err := ffjson.Unmarshal(b, &bit); err != nil {
					return nil, errors.Annotate(err, "unmarshal BitAssetData")
				}
				ret = append(ret, bit)

			default:
				util.DumpUnmarshaled(id.Type().String(), b)
				return nil, errors.Errorf("unable to parse GrapheneObject with ID %s", id)
			}
		}
	}

	return ret, nil
}

// CancelOrder cancels an order given by orderID
func (p *bitsharesAPI) CancelOrder(orderID types.GrapheneObject, broadcast bool) (*types.SignedTransaction, error) {
	defer p.SetDebug(false)

	resp, err := p.wsClient.CallAPI(0, "cancel_order", orderID.ID(), broadcast)
	if err != nil {
		return nil, err
	}

	p.Debug("cancel_order <", resp)

	ret := types.SignedTransaction{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Transaction")
	}

	return &ret, nil
}

//Debug prints the given object as human readable data to stdout
func (p bitsharesAPI) Debug(descr string, in interface{}) {
	if p.debug {
		util.Dump(descr, in)
	}
}

//SetDebug enables/disables function scoped debugging output.
//Note: due to readability reasons debug is switched off after return from call.
func (p *bitsharesAPI) SetDebug(debug bool) {
	p.debug = debug

	if p.wsClient != nil {
		p.wsClient.SetDebug(p.debug)
	}

	if p.rpcClient != nil {
		p.rpcClient.SetDebug(p.debug)
	}
}

func (p *bitsharesAPI) DatabaseAPIID() int {
	return p.databaseAPIID
}

func (p *bitsharesAPI) BroadcastAPIID() int {
	return p.broadcastAPIID
}

func (p *bitsharesAPI) HistoryAPIID() int {
	return p.historyAPIID
}

func (p *bitsharesAPI) CryptoAPIID() int {
	return p.cryptoAPIID
}

//CallWsAPI invokes a websocket API call
func (p *bitsharesAPI) CallWsAPI(apiID int, method string, args ...interface{}) (interface{}, error) {
	defer p.SetDebug(false)
	return p.wsClient.CallAPI(apiID, method, args...)
}

//OnError - hook your notify callback here
func (p *bitsharesAPI) OnNotify(subscriberID int, notifyFn func(msg interface{}) error) error {
	return p.wsClient.OnNotify(subscriberID, notifyFn)
}

//OnError - hook your error callback here
func (p *bitsharesAPI) OnError(errorFn client.ErrorFunc) {
	p.wsClient.OnError(errorFn)
}

//SetCredentials defines username and password for Websocket API login.
func (p *bitsharesAPI) SetCredentials(username, password string) {
	p.username = username
	p.password = password
}

// Connect initializes the API and connects underlying resources
// The websocket API is connected through the client provider.
func (p *bitsharesAPI) Connect() (err error) {

	if p.rpcClient != nil {
		if err := p.rpcClient.Connect(); err != nil {
			return errors.Annotate(err, "rpc connect")
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

	if err = config.SetCurrentConfig(chainID); err != nil {
		return errors.Annotate(err, "SetCurrentConfig")
	}

	return nil
}

func (p *bitsharesAPI) getAPIIDs() (err error) {
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

	p.cryptoAPIID, err = p.getAPIID("crypto")
	if err != nil {
		return errors.Annotate(err, "crypto")
	}

	return nil
}

//Close shuts the API down and closes underlying resources.
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
//wsEndpointURL: Is a mandatory websocket node URL.
//rpcEndpointURL: Is an optional RPC endpoint to your local `cli_wallet`.
//The use of wallet functions without this argument will throw an error.
//If you do not use wallet API, provide an empty string.
func New(wsEndpointURL, rpcEndpointURL string) BitsharesAPI {
	var rpcClient client.RPCClient
	if rpcEndpointURL != "" {
		rpcClient = client.NewRPCClient(rpcEndpointURL)
	}

	pr := NewSimpleClientProvider(wsEndpointURL)

	api := bitsharesAPI{
		wsClient:       pr,
		rpcClient:      rpcClient,
		databaseAPIID:  InvalidApiID,
		historyAPIID:   InvalidApiID,
		broadcastAPIID: InvalidApiID,
		cryptoAPIID:    InvalidApiID,
		debug:          false,
	}

	return &api
}

//NewWithAutoEndpoint creates a new BitsharesAPI interface with automatic node latency checking.
//It's best to use this API instance type for a long API lifecycle because the latency tester takes time to unleash its magic.
//startupEndpointURL: Iss a mandatory websocket node URL to startup the latency tester quickly.
//rpcEndpointURL: Is an optional RPC endpoint to your local `cli_wallet`.
//The use of wallet functions without this argument
//will throw an error. If you do not use wallet API, provide an empty string.
func NewWithAutoEndpoint(startupEndpointURL, rpcEndpointURL string) (BitsharesAPI, error) {
	var rpcClient client.RPCClient
	if rpcEndpointURL != "" {
		rpcClient = client.NewRPCClient(rpcEndpointURL)
	}

	api := &bitsharesAPI{
		rpcClient:      rpcClient,
		databaseAPIID:  InvalidApiID,
		historyAPIID:   InvalidApiID,
		broadcastAPIID: InvalidApiID,
		cryptoAPIID:    InvalidApiID,
		debug:          false,
	}

	pr, err := NewBestNodeClientProvider(startupEndpointURL, api)
	if err != nil {
		return nil, errors.Annotate(err, "NewBestNodeClientProvider")
	}

	api.wsClient = pr
	return api, nil
}
