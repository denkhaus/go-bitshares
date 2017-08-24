package tests

import (
	"testing"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/crypto"
	"github.com/denkhaus/bitshares/objects"
	"github.com/denkhaus/bitshares/operations"
	"github.com/stretchr/testify/suite"
)

type walletAPITest struct {
	suite.Suite
	TestAPI api.BitsharesAPI
}

func (suite *walletAPITest) SetupTest() {

	api := api.New(wsTestApiUrl)

	if err := api.Connect(); err != nil {
		suite.Fail(err.Error(), "Connect")
	}

	api.OnError(func(err error) {
		suite.Fail(err.Error(), "OnError")
	})

	suite.TestAPI = api
}

func (suite *walletAPITest) TearDown() {
	if err := suite.TestAPI.Close(); err != nil {
		suite.Fail(err.Error(), "Close")
	}
}

/*
func (suite *walletAPITest) Test_ListAssets() {
	res, err := suite.TestAPI.ListAssets("PEGF", 50)
	if err != nil {
		suite.Fail(err.Error(), "ListAssets")
	}

	suite.NotNil(res)
	suite.Len(res, 2)
	util.Dump("assets >", res)
} */

/* func (suite *walletAPITest) Test_GetBlock() {
	res, err := suite.TestAPI.GetBlock(10454132)
	if err != nil {
		suite.Fail(err.Error(), "GetBlock")
	}

	suite.NotNil(res)
	util.Dump("get_block >", res)
}
*/
/*
func (suite *walletAPITest) Test_Buy() {

	res, err := suite.TestAPI.Buy(AccountBuySell, AssetUSD, AssetBTS, 1111, 15, true)
	if err != nil {
		suite.Fail(err.Error(), "Buy")
	}

	util.Dump("buy <", res)
	suite.NotNil(res)
}
*/
/*
func (suite *walletAPITest) Test_GetAccountByName() {

	res, err := suite.TestAPI.GetAccountByName("denk-baum")
	if err != nil {
		suite.Fail(err.Error(), "GetAccountByName")
	}

	suite.NotNil(res)
	util.Dump("accounts >", res)
} */
/*
func (suite *walletAPITest) Test_GetLimitOrders() {

	res, err := suite.TestAPI.GetLimitOrders(AssetTEST, AssetBTS, 50)
	if err != nil {
		suite.Fail(err.Error(), "GetLimitOrders")
	}

	suite.NotNil(res)
	util.Dump("limitorders >", res)
}
*/

func (suite *walletAPITest) Test_CancelOrder() {

	op := operations.NewLimitOrderCancelOperation(
		*objects.NewGrapheneID("1.7.69314"),
	)
	op.FeePayingAccount = *TestAccount1ID

	priv, err := crypto.Decode(TestAccount1PrivKeyOwner)
	if err != nil {
		suite.Fail(err.Error(), "decode wif key")
	}

	_, err = suite.TestAPI.Broadcast([][]byte{priv}, AssetTEST, op)
	if err != nil {
		suite.Fail(err.Error(), "broadcast")
	}

}

/* func (suite *walletAPITest) Test_Transfer() {

	am := objects.AssetAmount{
		Amount: 1000,
		Asset:  *AssetTEST,
	}

	op := operations.TransferOperation{
		Extensions: []objects.Extension{},
		Amount:     am,
	}

	op.From.FromObjectID(TestAccount1ID.Id())
	op.To.FromObjectID(TestAccount2ID.Id())

	priv, err := crypto.Decode(TestAccount1PrivKey)
	if err != nil {
		suite.Fail(err.Error(), "decode wif key")
	}

	privKeys := [][]byte{priv}
	if err := suite.TestAPI.Broadcast(privKeys, op.Amount.Asset, &op); err != nil {
		suite.Fail(err.Error(), "broadcast")
	}

} */

func TestWalletApi(t *testing.T) {
	testSuite := new(walletAPITest)
	suite.Run(t, testSuite)
	testSuite.TearDown()
}
