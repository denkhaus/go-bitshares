package bitshares

import (
	"fmt"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/config"
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/denkhaus/logging"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type WalletAPI interface {
	Close() error
	Connect() error
	GetBlock(number uint64) (*types.Block, error)
	GetRelativeAccountHistory(account types.GrapheneObject, stop int64, limit int, start int64) (types.OperationRelativeHistories, error)
	GetDynamicGlobalProperties() (*types.DynamicGlobalProperties, error)
	BorrowAsset(account types.GrapheneObject, amountToBorrow string, symbolToBorrow types.GrapheneObject, amountOfCollateral string, broadcast bool) (*types.SignedTransaction, error)
	Buy(account types.GrapheneObject, base, quote types.GrapheneObject, rate string, amount string, broadcast bool) (*types.SignedTransaction, error)
	BuyEx(account types.GrapheneObject, base, quote types.GrapheneObject, rate float64, amount float64, broadcast bool) (*types.SignedTransaction, error)
	Info() (*types.Info, error)
	IsLocked() (bool, error)
	ListAccountBalances(account types.GrapheneObject) (types.AssetAmounts, error)
	Lock() error
	ReadMemo(memo *types.Memo) (string, error)
	Sell(account types.GrapheneObject, base, quote types.GrapheneObject, rate string, amount string, broadcast bool) (*types.SignedTransaction, error)
	SellEx(account types.GrapheneObject, base, quote types.GrapheneObject, rate float64, amount float64, broadcast bool) (*types.SignedTransaction, error)
	SellAsset(account types.GrapheneObject, amountToSell string, symbolToSell types.GrapheneObject, minToReceive string, symbolToReceive types.GrapheneObject, timeout uint32, fillOrKill bool, broadcast bool) (*types.SignedTransaction, error)
	SignTransaction(tx *types.SignedTransaction, broadcast bool) (*types.SignedTransaction, error)
	SerializeTransaction(tx *types.SignedTransaction) (string, error)
	//Transfer2(from, to types.GrapheneObject, amount string, asset types.GrapheneObject, memo string) (*types.SignedTransactionWithTransactionId, error)
	Unlock(password string) error
}

//NewWalletAPI creates a new WalletAPI interface.
//rpcEndpointURL: Is an RPC endpoint URL to your local `cli_wallet`.
func NewWalletAPI(rpcEndpointURL string) WalletAPI {
	api := &walletAPI{
		rpcClient: api.NewRPCClient(rpcEndpointURL),
	}

	return api
}

type walletAPI struct {
	rpcClient api.RPCClient
}

func (p *walletAPI) Connect() error {
	if err := p.rpcClient.Connect(); err != nil {
		return errors.Annotate(err, "Connect [rpc]")
	}

	info, err := p.Info()
	if err != nil {
		return errors.Annotate(err, "Info")
	}

	if err := config.SetCurrent(info.ChainID.String()); err != nil {
		return errors.Annotate(err, "SetCurrent")
	}

	return nil
}

//Close shuts the API down and closes underlying resources.
func (p *walletAPI) Close() error {
	if p.rpcClient != nil {
		if err := p.rpcClient.Close(); err != nil {
			return errors.Annotate(err, "Close [rpc]")
		}
		p.rpcClient = nil
	}

	return nil
}

// Lock locks the wallet
func (p *walletAPI) Lock() error {
	_, err := p.rpcClient.CallAPI("lock", types.EmptyParams)
	return err
}

// Unlock unlocks the wallet
func (p *walletAPI) Unlock(password string) error {
	_, err := p.rpcClient.CallAPI("unlock", password)
	return err
}

// IsLocked checks if wallet is locked.
func (p *walletAPI) IsLocked() (bool, error) {
	resp, err := p.rpcClient.CallAPI("is_locked", types.EmptyParams)

	if err != nil {
		return false, err
	}

	logging.DDumpJSON("is_locked <", resp)
	return resp.(bool), err
}

// Buy places a limit order attempting to buy one asset with another.
// This API call abstracts away some of the details of the sell_asset call to be more
// user friendly. All orders placed with buy never timeout and will not be killed if they
// cannot be filled immediately. If you wish for one of these parameters to be different,
// then sell_asset should be used instead.
//
// @param account The account buying the asset for another asset.
// @param base The name or id of the asset to buy.
// @param quote The name or id of the assest being offered as payment.
// @param rate The rate in base:quote at which you want to buy.
// @param amount the amount of base you want to buy.
// @param broadcast true to broadcast the transaction on the network.
// @returns The signed transaction selling the funds.
// @returns The error of operation.
func (p *walletAPI) Buy(account types.GrapheneObject, base, quote types.GrapheneObject, rate string, amount string, broadcast bool) (*types.SignedTransaction, error) {
	resp, err := p.rpcClient.CallAPI(
		"buy", account.ID(),
		base.ID(), quote.ID(),
		rate, amount, broadcast,
	)

	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("buy <", resp)

	ret := types.SignedTransaction{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Transaction")
	}

	return &ret, nil
}

// Sell places a limit order attempting to sell one asset for another.
// This API call abstracts away some of the details of the sell_asset call to be more
// user friendly. All orders placed with sell never timeout and will not be killed if they
// cannot be filled immediately. If you wish for one of these parameters to be different,
// then sell_asset should be used instead.
//
// @param account the account providing the asset being sold, and which will receive the processed of the sale.
// @param base The name or id of the asset to sell.
// @param quote The name or id of the asset to receive.
// @param rate The rate in base:quote at which you want to sell.
// @param amount The amount of base you want to sell.
// @param broadcast true to broadcast the transaction on the network.
// @returns The signed transaction selling the funds.
// @returns The error of operation.
func (p *walletAPI) Sell(account types.GrapheneObject, base, quote types.GrapheneObject, rate string, amount string, broadcast bool) (*types.SignedTransaction, error) {
	resp, err := p.rpcClient.CallAPI(
		"sell", account.ID(),
		base.ID(), quote.ID(),
		rate, amount, broadcast,
	)

	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("sell <", resp)

	ret := types.SignedTransaction{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Transaction")
	}

	return &ret, nil
}

func (p *walletAPI) BuyEx(account types.GrapheneObject, base, quote types.GrapheneObject,
	rate float64, amount float64, broadcast bool) (*types.SignedTransaction, error) {
	//TODO: use proper precision, avoid rounding
	minToReceive := fmt.Sprintf("%f", amount)
	amountToSell := fmt.Sprintf("%f", rate*amount)

	return p.SellAsset(account, amountToSell, quote, minToReceive, base, 0, false, broadcast)
}

func (p *walletAPI) SellEx(account types.GrapheneObject, base, quote types.GrapheneObject,
	rate float64, amount float64, broadcast bool) (*types.SignedTransaction, error) {

	//TODO: use proper precision, avoid rounding
	amountToSell := fmt.Sprintf("%f", amount)
	minToReceive := fmt.Sprintf("%f", rate*amount)

	return p.SellAsset(account, amountToSell, base, minToReceive, quote, 0, false, broadcast)
}

// SellAsset
func (p *walletAPI) SellAsset(account types.GrapheneObject, amountToSell string, symbolToSell types.GrapheneObject,
	minToReceive string, symbolToReceive types.GrapheneObject, timeout uint32, fillOrKill bool, broadcast bool) (*types.SignedTransaction, error) {
	resp, err := p.rpcClient.CallAPI(
		"sell_asset", account.ID(),
		amountToSell, symbolToSell.ID(),
		minToReceive, symbolToReceive.ID(),
		timeout, fillOrKill, broadcast,
	)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("sell_asset <", resp)

	ret := types.SignedTransaction{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Transaction")
	}

	return &ret, nil
}

// BorrowAsset borrows an asset or update the debt/collateral ratio for the loan.
// @param account: the id of the account associated with the transaction.
// @param amountToBorrow: the amount of the asset being borrowed. Make this value negative to pay back debt.
// @param symbolToBorrow: the symbol or id of the asset being borrowed.
// @param amountOfCollateral: the amount of the backing asset to add to your collateral position. Make this negative to claim back some of your collateral. The backing asset is defined in the bitasset_options for the asset being borrowed.
// @param broadcast: true to broadcast the transaction on the network
func (p *walletAPI) BorrowAsset(account types.GrapheneObject, amountToBorrow string, symbolToBorrow types.GrapheneObject,
	amountOfCollateral string, broadcast bool) (*types.SignedTransaction, error) {

	resp, err := p.rpcClient.CallAPI(
		"borrow_asset", account.ID(),
		amountToBorrow, symbolToBorrow.ID(),
		amountOfCollateral, broadcast,
	)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("borrow_asset <", resp)

	ret := types.SignedTransaction{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Transaction")
	}

	return &ret, nil
}

func (p *walletAPI) ListAccountBalances(account types.GrapheneObject) (types.AssetAmounts, error) {
	resp, err := p.rpcClient.CallAPI("list_account_balances", account.ID())
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("list_account_balances <", resp)

	ret := types.AssetAmounts{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal AssetAmounts")
	}

	return ret, nil
}

// SerializeTransaction converts a signed_transaction in JSON form to its binary representation.
// @param tx the transaction to serialize
// Returns the binary form of the transaction. It will not be hex encoded, this returns a raw string that may have null characters embedded in it.
func (p *walletAPI) SerializeTransaction(tx *types.SignedTransaction) (string, error) {
	resp, err := p.rpcClient.CallAPI("serialize_transaction", tx)
	if err != nil {
		return "", err
	}

	logging.DDumpJSON("serialize_transaction <", resp)

	return resp.(string), nil
}

// SignTransaction signs a transaction
// @param tx the transaction to sign
// @param broadcast bool defines if the transaction should be broadcasted
// Returns the signed transaction.
func (p *walletAPI) SignTransaction(tx *types.SignedTransaction, broadcast bool) (*types.SignedTransaction, error) {
	resp, err := p.rpcClient.CallAPI("sign_transaction", tx, broadcast)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("sign_transaction <", resp)

	ret := types.SignedTransaction{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Transaction")
	}

	return &ret, nil

}

func (p *walletAPI) ReadMemo(memo *types.Memo) (string, error) {
	resp, err := p.rpcClient.CallAPI("read_memo", memo)
	if err != nil {
		return "", err
	}
	if msg, ok := resp.(string); ok {
		return msg, nil
	}
	return "", nil
}

//GetBlock retrieves a block by number
func (p *walletAPI) GetBlock(number uint64) (*types.Block, error) {
	resp, err := p.rpcClient.CallAPI("get_block", number)
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

//GetRelativeAccountHistory gets operations relevant to the specified account referenced by an event numbering specific to the account. The current number of operations for the account can be found in the account statistics (or use 0 for start).
//
//Parameters
//   account_id_or_name	The account ID or name whose history should be queried
//   stop	Sequence number of earliest operation. 0 is default and will query 'limit' number of operations.
//   limit	Maximum number of operations to retrieve (must not exceed 100)
//   start	Sequence number of the most recent operation to retrieve. 0 is default, which will start querying from the most recent operation.
//
//Returns
//   A list of operations performed by account, ordered from most recent to oldest.
func (p *walletAPI) GetRelativeAccountHistory(account types.GrapheneObject, stop int64, limit int, start int64) (types.OperationRelativeHistories, error) {
	if limit > GetAccountHistoryLimit {
		limit = GetAccountHistoryLimit
	}

	resp, err := p.rpcClient.CallAPI("get_relative_account_history", account.ID(), stop, limit, start)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("get_relative_account_history <", resp)

	ret := types.OperationRelativeHistories{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Histories")
	}

	return ret, nil
}

// GetDynamicGlobalProperties returns the block chain’s rapidly-changing properties.
// The returned object contains information that changes every block
// interval such as the head block number, the next witness, etc.
//
// See
//   get_global_properties() for less-frequently changing properties
// Return
//   the dynamic global properties
func (p *walletAPI) GetDynamicGlobalProperties() (*types.DynamicGlobalProperties, error) {
	resp, err := p.rpcClient.CallAPI("get_dynamic_global_properties", types.EmptyParams)
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

func (p *walletAPI) Info() (*types.Info, error) {
	resp, err := p.rpcClient.CallAPI("info", types.EmptyParams)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("info <", resp)

	var ret types.Info
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal info")
	}

	return &ret, nil
}

//Transfer2 works just like transfer, except it always broadcasts and
//returns the transaction ID along with the signed transaction.
// func (p *walletAPI) Transfer2(from, to types.GrapheneObject, amount string, asset types.GrapheneObject, memo string) (*types.SignedTransactionWithTransactionId, error) {
// 	if p.rpcClient == nil {
// 		return nil, types.ErrRPCClientNotInitialized
// 	}
// 	resp, err := p.rpcClient.CallAPI("transfer2", from.ID(), to.ID(), amount, asset.ID(), memo)
// 	if err != nil {
// 		return nil, err
// 	}
// 	ret := types.SignedTransactionWithTransactionId{}
// 	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
// 		return nil, errors.Annotate(err, "unmarshal Transaction")
// 	}
// 	return &ret, nil
// }
