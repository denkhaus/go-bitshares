package api

import (
	"fmt"

	"github.com/denkhaus/bitshares/objects"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

// lock the wallet
func (p *bitsharesAPI) WalletLock() error {
	_, err := p.rpcClient.CallAPI("lock", EmptyParams)
	return err
}

// unlock the wallet
func (p *bitsharesAPI) WalletUnlock(password string) error {
	_, err := p.rpcClient.CallAPI("unlock", password)

	return err
}

// Check if wallet is locked.
func (p *bitsharesAPI) WalletIsLocked() (bool, error) {
	resp, err := p.rpcClient.CallAPI("is_locked", EmptyParams)

	return resp.(bool), err
}

// Place a limit order attempting to buy one asset with another.
//
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
// @param The signed transaction selling the funds.
// @param error of operation.
func (p *bitsharesAPI) Buy(account objects.GrapheneObject, base, quote objects.GrapheneObject,
	rate string, amount string, broadcast bool) (*objects.Transaction, error) {

	resp, err := p.rpcClient.CallAPI("buy", account.Id(), base.Id(), quote.Id(), rate, amount, broadcast)
	if err != nil {
		return nil, err
	}

	util.Dump("buy >", resp)

	ret := objects.Transaction{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Transaction")
	}

	return &ret, nil
}

//Place a limit order attempting to sell one asset for another.
//
// This API call abstracts away some of the details of the sell_asset call to be more
// user friendly. All orders placed with sell never timeout and will not be killed if they
// cannot be filled immediately. If you wish for one of these parameters to be different,
// then sell_asset should be used instead.
//
// @param account the account providing the asset being sold, and which will receive the processed of the sale.
// @param base The name or id of the asset to sell.
// @param quote The name or id of the asset to recieve.
// @param rate The rate in base:quote at which you want to sell.
// @param amount The amount of base you want to sell.
// @param broadcast true to broadcast the transaction on the network.
// @returns The signed transaction selling the funds.
// @param error of operation.
func (p *bitsharesAPI) Sell(account objects.GrapheneObject, base, quote objects.GrapheneObject,
	rate string, amount string, broadcast bool) (*objects.Transaction, error) {
	resp, err := p.rpcClient.CallAPI("sell", account.Id(), base.Id(), quote.Id(), rate, amount, broadcast)
	if err != nil {
		return nil, err
	}

	ret := objects.Transaction{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Transaction")
	}

	return &ret, nil
}

func (p *bitsharesAPI) BuyEx(account objects.GrapheneObject, base, quote objects.GrapheneObject,
	rate float64, amount float64, broadcast bool) (*objects.Transaction, error) {

	//TODO: use proper precision, avoid rounding
	minToReceive := fmt.Sprintf("%f", amount)
	amountToSell := fmt.Sprintf("%f", rate*amount)

	return p.SellAsset(account, amountToSell, quote, minToReceive, base, 0, false, broadcast)
}

func (p *bitsharesAPI) SellEx(account objects.GrapheneObject, base, quote objects.GrapheneObject,
	rate float64, amount float64, broadcast bool) (*objects.Transaction, error) {

	//TODO: use proper precision, avoid rounding
	amountToSell := fmt.Sprintf("%f", amount)
	minToReceive := fmt.Sprintf("%f", rate*amount)

	return p.SellAsset(account, amountToSell, base, minToReceive, quote, 0, false, broadcast)
}

func (p *bitsharesAPI) SellAsset(account objects.GrapheneObject,
	amountToSell string, symbolToSell objects.GrapheneObject,
	minToReceive string, symbolToReceive objects.GrapheneObject,
	timeout uint32, fillOrKill bool, broadcast bool) (*objects.Transaction, error) {

	resp, err := p.rpcClient.CallAPI("sell_asset", account.Id(),
		amountToSell, symbolToSell.Id(),
		minToReceive, symbolToReceive.Id(),
		timeout, fillOrKill, broadcast,
	)
	if err != nil {
		return nil, err
	}

	//util.Dump("sell_asset >", resp)

	ret := objects.Transaction{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Transaction")
	}

	return &ret, nil
}

//Borrow an asset or update the debt/collateral ratio for the loan.
// @param account: the id of the account associated with the transaction.
// @param amountToBorrow: the amount of the asset being borrowed. Make this value negative to pay back debt.
// @param symbolToBorrow: the symbol or id of the asset being borrowed.
// @param amountOfCollateral: the amount of the backing asset to add to your collateral position. Make this negative to claim back some of your collateral. The backing asset is defined in the bitasset_options for the asset being borrowed.
// @param broadcast: true to broadcast the transaction on the network
func (p *bitsharesAPI) BorrowAsset(account objects.GrapheneObject,
	amountToBorrow string, symbolToBorrow objects.GrapheneObject,
	amountOfCollateral string, broadcast bool) (*objects.Transaction, error) {

	resp, err := p.rpcClient.CallAPI("borrow_asset", account.Id(),
		amountToBorrow, symbolToBorrow.Id(),
		amountOfCollateral, broadcast,
	)
	if err != nil {
		return nil, err
	}

	//util.Dump("borrow_asset >", resp)

	ret := objects.Transaction{}
	if err := ffjson.Unmarshal(util.ToBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Transaction")
	}

	return &ret, nil
}

func (p *bitsharesAPI) ListAccountBalances(account objects.GrapheneObject) ([]objects.AssetAmount, error) {

	resp, err := p.rpcClient.CallAPI("list_account_balances", account.Id())
	if err != nil {
		return nil, err
	}

	//util.Dump("list_account_balances >", resp)

	data := resp.([]interface{})
	ret := make([]objects.AssetAmount, len(data))

	for idx, a := range data {
		if err := ffjson.Unmarshal(util.ToBytes(a), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal AssetAmount")
		}
	}

	return ret, nil
}
