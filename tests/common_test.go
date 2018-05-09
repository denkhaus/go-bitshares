package tests

import (
	"testing"
	"time"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/objects"
	"github.com/stretchr/testify/suite"
)

type commonTest struct {
	suite.Suite
	TestAPI api.BitsharesAPI
}

func (suite *commonTest) SetupTest() {
	api := api.New(wsFullApiUrl, rpcApiUrl)
	if err := api.Connect(); err != nil {
		suite.Fail(err.Error(), "Connect")
	}

	api.OnError(func(err error) {
		suite.Fail(err.Error(), "OnError")
	})

	suite.TestAPI = api
}

func (suite *commonTest) TearDown() {
	if err := suite.TestAPI.Close(); err != nil {
		suite.Fail(err.Error(), "Close")
	}
}

func (suite *commonTest) Test_GetChainID() {

	res, err := suite.TestAPI.GetChainID()
	if err != nil {
		suite.Fail(err.Error(), "GetChainID")
	}

	suite.Equal(res, ChainIDBitSharesFull)
}

func (suite *commonTest) Test_GetAccountBalances() {

	res, err := suite.TestAPI.GetAccountBalances(UserID2, AssetBTS)
	if err != nil {
		suite.Fail(err.Error(), "GetAccountBalances 1")
	}

	suite.NotNil(res)
	//util.Dump("balance bts >", res)

	res, err = suite.TestAPI.GetAccountBalances(UserID2)
	if err != nil {
		suite.Fail(err.Error(), "GetAccountBalances 2")
	}

	suite.NotNil(res)
	//util.Dump("balances all >", res)
}

func (suite *commonTest) Test_GetAccounts() {

	res, err := suite.TestAPI.GetAccounts(UserID3)
	if err != nil {
		suite.Fail(err.Error(), "GetAccounts")
	}

	suite.NotNil(res)
	suite.Len(res, 1)
}

func (suite *commonTest) Test_GetObjects() {

	res, err := suite.TestAPI.GetObjects(
		UserID1,
		AssetCNY,
		BitAssetDataCNY,
		LimitOrder1,
		CallOrder1,
		SettleOrder1,
	)

	if err != nil {
		suite.Fail(err.Error(), "GetObjects")
	}

	suite.NotNil(res)
	suite.Len(res, 6)
	//util.Dump("objects >", res)
}

func (suite *commonTest) Test_GetBlock() {
	res, err := suite.TestAPI.GetBlock(1316)
	if err != nil {
		suite.Fail(err.Error(), "GetBlock")
	}

	suite.NotNil(res)
	//util.Dump("get_block >", res)
}

func (suite *commonTest) Test_GetDynamicGlobalProperties() {
	res, err := suite.TestAPI.GetDynamicGlobalProperties()
	if err != nil {
		suite.Fail(err.Error(), "GetDynamicGlobalProperties")
	}

	suite.NotNil(res)
	//util.Dump("dynamic global properties >", res)
}

func (suite *commonTest) Test_GetAccountByName() {
	res, err := suite.TestAPI.GetAccountByName("openledger")
	if err != nil {
		suite.Fail(err.Error(), "GetAccountByName")
	}

	suite.NotNil(res)
	//util.Dump("accounts >", res)
}

func (suite *commonTest) Test_GetTradeHistory() {
	dtTo := time.Now().UTC()

	dtFrom := dtTo.Add(-time.Hour * 24)
	res, err := suite.TestAPI.GetTradeHistory(AssetBTS, AssetHERO, dtTo, dtFrom, 50)

	if err != nil {
		suite.Fail(err.Error(), "GetTradeHistory")
	}

	suite.NotNil(res)
	//util.Dump("markettrades >", res)
}

func (suite *commonTest) Test_GetLimitOrders() {

	res, err := suite.TestAPI.GetLimitOrders(AssetCNY, AssetBTS, 50)
	if err != nil {
		suite.Fail(err.Error(), "GetLimitOrders")
	}

	suite.NotNil(res)
	//util.Dump("limitorders >", res)
}

func (suite *commonTest) Test_GetCallOrders() {
	res, err := suite.TestAPI.GetCallOrders(AssetUSD, 50)
	if err != nil {
		suite.Fail(err.Error(), "GetCallOrders")
	}

	suite.NotNil(res)
	//	util.Dump("callorders >", res)
}

func (suite *commonTest) Test_GetMarginPositions() {
	res, err := suite.TestAPI.GetMarginPositions(UserID2)
	if err != nil {
		suite.Fail(err.Error(), "GetMarginPositions")
	}

	suite.NotNil(res)
	//util.Dump("marginpositions >", res)
}

func (suite *commonTest) Test_GetSettleOrders() {
	res, err := suite.TestAPI.GetSettleOrders(AssetCNY, 50)
	if err != nil {
		suite.Fail(err.Error(), "GetSettleOrders")
	}

	suite.NotNil(res)
	//util.Dump("settleorders >", res)
}

func (suite *commonTest) Test_ListAssets() {
	res, err := suite.TestAPI.ListAssets("OPEN.DASH", 2)
	if err != nil {
		suite.Fail(err.Error(), "ListAssets")
	}

	suite.NotNil(res)
	suite.Len(res, 2)
	//util.Dump("assets >", res)
}

func (suite *commonTest) Test_GetAccountHistory() {

	user := objects.NewGrapheneID("1.2.96393")
	start := objects.NewGrapheneID("1.11.187698971")
	stop := objects.NewGrapheneID("1.11.187658388")

	res, err := suite.TestAPI.GetAccountHistory(user, stop, 30, start)
	if err != nil {
		suite.Fail(err.Error(), "GetAccountHistory")
	}

	suite.NotNil(res)
	//util.Dump("history >", res)
}

func TestCommon(t *testing.T) {
	testSuite := new(commonTest)
	suite.Run(t, testSuite)
	testSuite.TearDown()
}
