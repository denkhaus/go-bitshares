package tests

import (
	"testing"

	"github.com/denkhaus/bitshares/api"
	"github.com/stretchr/testify/suite"
)

type walletAPITest struct {
	suite.Suite
	TestAPI api.BitsharesAPI
}

func (suite *walletAPITest) SetupTest() {

	api := api.New(WsTestApiUrl, RpcApiUrl)

	if err := api.Connect(); err != nil {
		suite.FailNow(err.Error(), "Connect")
	}

	api.OnError(func(err error) {
		suite.FailNow(err.Error(), "OnError")
	})

	suite.TestAPI = api
}

func (suite *walletAPITest) TearDown() {
	if err := suite.TestAPI.Close(); err != nil {
		suite.FailNow(err.Error(), "Close")
	}
}

// func (suite *walletAPITest) Test_ListAssets() {
// 	res, err := suite.TestAPI.ListAssets("PEG.FAKEUSD", 2)
// 	if err != nil {
// 		suite.FailNow(err.Error(), "ListAssets")
// 	}

// 	suite.NotNil(res)
// 	suite.Len(res, 2)
// 	util.Dump("assets >", res)
// }

/* func (suite *walletAPITest) Test_GetBlock() {
	res, err := suite.TestAPI.GetBlock(10454132)
	if err != nil {
		suite.FailNow(err.Error(), "GetBlock")
	}

	suite.NotNil(res)
	util.Dump("get_block >", res)
} */

func (suite *walletAPITest) Test_ChainConfig() {
	res, err := suite.TestAPI.GetChainID()
	if err != nil {
		suite.FailNow(err.Error(), "GetChainID")
	}

	suite.Equal(ChainIDBitSharesTest, res)
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
/* func (suite *walletAPITest) Test_Transfer() {

	am := types.AssetAmount{
		Amount: 1000,
		Asset:  *AssetTEST,
	}

	op := operations.TransferOperation{
		Extensions: []types.Extension{},
		Amount:     am,
	}

	op.From.FromObjectID(TestAccount1ID.Id())
	op.To.FromObjectID(TestAccount2ID.Id())

	priv, err := crypto.Decode(TestAccount1PrivKey)
	if err != nil {
		suite.FailNow(err.Error(), "decode wif key")
	}

	privKeys := [][]byte{priv}
	if err := suite.TestAPI.Broadcast(privKeys, op.Amount.Asset, &op); err != nil {
		suite.FailNow(err.Error(), "broadcast")
	}

} */

func TestWalletApi(t *testing.T) {
	testSuite := new(walletAPITest)
	suite.Run(t, testSuite)
	testSuite.TearDown()
}
