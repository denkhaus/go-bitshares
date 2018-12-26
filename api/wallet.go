package api

import (
	"fmt"

	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/denkhaus/logging"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

// WalletLock locks the wallet
func (p *bitsharesAPI) WalletLock() error {
	if p.rpcClient == nil {
		return types.ErrRPCClientNotInitialized
	}

	_, err := p.rpcClient.CallAPI("lock", types.EmptyParams)
	return err
}

// WalletUnlock unlocks the wallet
func (p *bitsharesAPI) WalletUnlock(password string) error {
	if p.rpcClient == nil {
		return types.ErrRPCClientNotInitialized
	}

	_, err := p.rpcClient.CallAPI("unlock", password)
	return err
}

// WalletIsLocked checks if wallet is locked.
func (p *bitsharesAPI) WalletIsLocked() (bool, error) {
	if p.rpcClient == nil {
		return false, types.ErrRPCClientNotInitialized
	}

	resp, err := p.rpcClient.CallAPI("is_locked", types.EmptyParams)

	if err != nil {
		return false, err
	}

	logging.DDumpJSON("is_locked <", resp)

	return resp.(bool), err
}

// WalletBuy places a limit order attempting to buy one asset with another.
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
func (p *bitsharesAPI) WalletBuy(account types.GrapheneObject, base, quote types.GrapheneObject, rate string, amount string, broadcast bool) (*types.SignedTransaction, error) {
	if p.rpcClient == nil {
		return nil, types.ErrRPCClientNotInitialized
	}

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

// WalletSell places a limit order attempting to sell one asset for another.
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
func (p *bitsharesAPI) WalletSell(account types.GrapheneObject, base, quote types.GrapheneObject, rate string, amount string, broadcast bool) (*types.SignedTransaction, error) {
	if p.rpcClient == nil {
		return nil, types.ErrRPCClientNotInitialized
	}

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

func (p *bitsharesAPI) WalletBuyEx(account types.GrapheneObject, base, quote types.GrapheneObject,
	rate float64, amount float64, broadcast bool) (*types.SignedTransaction, error) {
	//TODO: use proper precision, avoid rounding
	minToReceive := fmt.Sprintf("%f", amount)
	amountToSell := fmt.Sprintf("%f", rate*amount)

	return p.WalletSellAsset(account, amountToSell, quote, minToReceive, base, 0, false, broadcast)
}

func (p *bitsharesAPI) WalletSellEx(account types.GrapheneObject, base, quote types.GrapheneObject,
	rate float64, amount float64, broadcast bool) (*types.SignedTransaction, error) {

	//TODO: use proper precision, avoid rounding
	amountToSell := fmt.Sprintf("%f", amount)
	minToReceive := fmt.Sprintf("%f", rate*amount)

	return p.WalletSellAsset(account, amountToSell, base, minToReceive, quote, 0, false, broadcast)
}

// SellAsset
func (p *bitsharesAPI) WalletSellAsset(account types.GrapheneObject, amountToSell string, symbolToSell types.GrapheneObject,
	minToReceive string, symbolToReceive types.GrapheneObject, timeout uint32, fillOrKill bool, broadcast bool) (*types.SignedTransaction, error) {
	if p.rpcClient == nil {
		return nil, types.ErrRPCClientNotInitialized
	}

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

// WalletBorrowAsset borrows an asset or update the debt/collateral ratio for the loan.
// @param account: the id of the account associated with the transaction.
// @param amountToBorrow: the amount of the asset being borrowed. Make this value negative to pay back debt.
// @param symbolToBorrow: the symbol or id of the asset being borrowed.
// @param amountOfCollateral: the amount of the backing asset to add to your collateral position. Make this negative to claim back some of your collateral. The backing asset is defined in the bitasset_options for the asset being borrowed.
// @param broadcast: true to broadcast the transaction on the network
func (p *bitsharesAPI) WalletBorrowAsset(account types.GrapheneObject, amountToBorrow string, symbolToBorrow types.GrapheneObject,
	amountOfCollateral string, broadcast bool) (*types.SignedTransaction, error) {

	if p.rpcClient == nil {
		return nil, types.ErrRPCClientNotInitialized
	}

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

func (p *bitsharesAPI) WalletListAccountBalances(account types.GrapheneObject) (types.AssetAmounts, error) {

	if p.rpcClient == nil {
		return nil, types.ErrRPCClientNotInitialized
	}

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
func (p *bitsharesAPI) WalletSerializeTransaction(tx *types.SignedTransaction) (string, error) {
	if p.rpcClient == nil {
		return "", types.ErrRPCClientNotInitialized
	}

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
func (p *bitsharesAPI) WalletSignTransaction(tx *types.SignedTransaction, broadcast bool) (*types.SignedTransaction, error) {
	if p.rpcClient == nil {
		return nil, types.ErrRPCClientNotInitialized
	}

	resp, err := p.rpcClient.CallAPI("sign_transaction", tx, broadcast)
	if err != nil {
		return nil, err
	}

	logging.DDumpJSON("wallet sign_transaction <", resp)

	ret := types.SignedTransaction{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Transaction")
	}

	return &ret, nil

}

func (p *bitsharesAPI) WalletReadMemo(memo *types.Memo) (string, error) {
	if p.rpcClient == nil {
		return "", types.ErrRPCClientNotInitialized
	}
	resp, err := p.rpcClient.CallAPI("read_memo", memo)
	if err != nil {
		return "", err
	}
	if msg, ok := resp.(string); ok {
		return msg, nil
	}
	return "", nil
}

//WalletTransfer2 works just like transfer, except it always broadcasts and returns the transaction ID along with the signed transaction.
// func (p *bitsharesAPI) WalletTransfer2(from, to types.GrapheneObject, amount string, asset types.GrapheneObject, memo string) (*types.SignedTransactionWithTransactionId, error) {
// 	if p.rpcClient == nil {
// 		return nil, types.ErrRPCClientNotInitialized
// 	}
// 	resp, err := p.rpcClient.CallAPI("transfer2", from, to, amount, asset, memo)
// 	if err != nil {
// 		return nil, err
// 	}
// 	ret := types.SignedTransactionWithTransactionId{}
// 	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
// 		return nil, errors.Annotate(err, "unmarshal Transaction")
// 	}
// 	return &ret, nil
// }

//WalletGetBlock retrieves a block by number
func (p *bitsharesAPI) WalletGetBlock(number uint64) (*types.Block, error) {
	if p.rpcClient == nil {
		return nil, types.ErrRPCClientNotInitialized
	}

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

//WalletGetRelativeAccountHistory gets operations relevant to the specified account referenced by an event numbering specific to the account. The current number of operations for the account can be found in the account statistics (or use 0 for start).
//
//Parameters
//   account_id_or_name	The account ID or name whose history should be queried
//   stop	Sequence number of earliest operation. 0 is default and will query 'limit' number of operations.
//   limit	Maximum number of operations to retrieve (must not exceed 100)
//   start	Sequence number of the most recent operation to retrieve. 0 is default, which will start querying from the most recent operation.
//
//Returns
//   A list of operations performed by account, ordered from most recent to oldest.
func (p *bitsharesAPI) WalletGetRelativeAccountHistory(account types.GrapheneObject, stop int64, limit int, start int64) (types.OperationRelativeHistories, error) {
	if p.rpcClient == nil {
		return nil, types.ErrRPCClientNotInitialized
	}

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

// WalletGetDynamicGlobalProperties returns the block chain’s rapidly-changing properties.
// The returned object contains information that changes every block
// interval such as the head block number, the next witness, etc.
//
// See
//   get_global_properties() for less-frequently changing properties
// Return
//   the dynamic global properties
func (p *bitsharesAPI) WalletGetDynamicGlobalProperties() (*types.DynamicGlobalProperties, error) {
	if p.rpcClient == nil {
		return nil, types.ErrRPCClientNotInitialized
	}

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
