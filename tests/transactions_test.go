package tests

import (
	"testing"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/crypto"
	"github.com/denkhaus/bitshares/operations"
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/stretchr/testify/suite"
)

type walletAPITest struct {
	suite.Suite
	TestAPI api.BitsharesAPI
	KeyBag  *crypto.KeyBag
}

func (suite *walletAPITest) SetupTest() {
	suite.TestAPI = NewTestAPI(suite.T(), WsTestApiUrl)
	suite.KeyBag = crypto.NewKeyBag()

	if err := suite.KeyBag.Add(TestAccount1PrivKeyActive); err != nil {
		suite.FailNow(err.Error(), "KeyBag.Add 1")
	}
	if err := suite.KeyBag.Add(TestAccount3PrivKeyActive); err != nil {
		suite.FailNow(err.Error(), "KeyBag.Add 2")
	}
}

func (suite *walletAPITest) TearDown() {
	if err := suite.TestAPI.Close(); err != nil {
		suite.FailNow(err.Error(), "Close")
	}
}

func (suite *walletAPITest) Test_ChainConfig() {
	res, err := suite.TestAPI.GetChainID()
	if err != nil {
		suite.FailNow(err.Error(), "GetChainID")
	}

	suite.Equal(ChainIDBitSharesTest, res)
}

func (suite *walletAPITest) Test_SignTransaction() {
	op := operations.TransferOperation{
		Extensions: types.Extensions{},
		Amount: types.AssetAmount{
			Amount: 1000,
			Asset:  *AssetTEST,
		},
		From: *TestAccount1ID,
		To:   *TestAccount2ID,
	}

	trx, err := suite.TestAPI.SignTransaction(suite.KeyBag, AssetTEST, &op)
	if err != nil {
		suite.FailNow(err.Error(), "SignTransaction")
	}

	util.Dump("signed trx <", trx)
}

/*
func (suite *walletAPITest) Test_Buy() {

	res, err := suite.TestAPI.Buy(AccountBuySell, AssetUSD, AssetBTS, 1111, 15, true)
	if err != nil {
		suite.FailNow(err.Error(), "Buy")
	}

	util.Dump("buy <", res)
	suite.NotNil(res)
}
*/
/*
func (suite *walletAPITest) Test_GetAccountByName() {

	res, err := suite.TestAPI.GetAccountByName("denk-haus")
	if err != nil {
		suite.FailNow(err.Error(), "GetAccountByName")
	}

	suite.NotNil(res)
	util.Dump("accounts >", res)
} */

/* func (suite *walletAPITest) Test_GetLimitOrders() {

	res, err := suite.TestAPI.GetLimitOrders(AssetTEST, AssetPEGFAKEUSD, 50)
	if err != nil {
		suite.FailNow(err.Error(), "GetLimitOrders")
	}

	suite.NotNil(res)
	util.Dump("limitorders >", res)
}


func (suite *walletAPITest) Test_CancelOrder() {

	op := types.NewLimitOrderCancelOperation(
		*types.NewGrapheneID("1.7.69314"),
	)
	op.FeePayingAccount = *TestAccount1ID

	_, err := suite.TestAPI.Broadcast([]string{TestAccount1PrivKeyActive}, AssetTEST, op)
	if err != nil {
		suite.FailNow(err.Error(), "broadcast")
	}

}
*/

func TestWalletApi(t *testing.T) {
	testSuite := new(walletAPITest)
	suite.Run(t, testSuite)
	testSuite.TearDown()
}
