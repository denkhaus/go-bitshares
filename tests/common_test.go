package tests

import (
	"testing"
	"time"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/types"
	"github.com/stretchr/testify/suite"

	//import operations to initialize types.OperationMap
	_ "github.com/denkhaus/bitshares/operations"
)

type commonTest struct {
	suite.Suite
	TestAPI api.BitsharesAPI
}

func (suite *commonTest) SetupTest() {
	suite.TestAPI = NewTestAPI(suite.T(), WsFullApiUrl)
}

func (suite *commonTest) TearDownTest() {
	if err := suite.TestAPI.Close(); err != nil {
		suite.FailNow(err.Error(), "Close")
	}
}

func (suite *commonTest) Test_GetChainID() {
	res, err := suite.TestAPI.GetChainID()
	if err != nil {
		suite.FailNow(err.Error(), "GetChainID")
	}

	suite.Equal(res, ChainIDBitSharesFull)
}

func (suite *commonTest) Test_GetAccountBalances() {
	res, err := suite.TestAPI.GetAccountBalances(UserID2, AssetBTS)
	if err != nil {
		suite.FailNow(err.Error(), "GetAccountBalances 1")
	}

	suite.NotNil(res)
	//util.Dump("balance bts >", res)

	res, err = suite.TestAPI.GetAccountBalances(UserID2)
	if err != nil {
		suite.FailNow(err.Error(), "GetAccountBalances 2")
	}

	suite.NotNil(res)
	//util.Dump("balances all >", res)
}

func (suite *commonTest) Test_GetAccounts() {
	suite.TestAPI.SetDebug(false)
	res, err := suite.TestAPI.GetAccounts(UserID2) //, UserID3, UserID4)
	if err != nil {
		suite.FailNow(err.Error(), "GetAccounts")
	}

	suite.NotNil(res)
	suite.Len(res, 1)

	//util.Dump("get accounts >", res)
}

func (suite *commonTest) Test_GetObjects() {
	res, err := suite.TestAPI.GetObjects(
		UserID1,
		AssetCNY,
		BitAssetDataCNY,
		LimitOrder1,
		CallOrder1,
		SettleOrder1,
		OperationHistory1,
	)

	if err != nil {
		suite.FailNow(err.Error(), "GetObjects")
	}

	suite.NotNil(res)
	suite.Len(res, 7)
	//util.Dump("objects >", res)
}

func (suite *commonTest) Test_GetBlock() {
	res, err := suite.TestAPI.GetBlock(26867161)
	if err != nil {
		suite.FailNow(err.Error(), "GetBlock")
	}

	suite.NotNil(res)
	//util.Dump("get_block >", res)
}

func (suite *commonTest) Test_GetDynamicGlobalProperties() {
	res, err := suite.TestAPI.GetDynamicGlobalProperties()
	if err != nil {
		suite.FailNow(err.Error(), "GetDynamicGlobalProperties")
	}

	suite.NotNil(res)
	//util.Dump("dynamic global properties >", res)
}

func (suite *commonTest) Test_GetAccountByName() {
	res, err := suite.TestAPI.GetAccountByName("openledger")
	if err != nil {
		suite.FailNow(err.Error(), "GetAccountByName")
	}

	suite.NotNil(res)
	//util.Dump("accounts >", res)
}

func (suite *commonTest) Test_GetTradeHistory() {
	dtTo := time.Now().UTC()

	dtFrom := dtTo.Add(-time.Hour * 24)
	res, err := suite.TestAPI.GetTradeHistory(AssetBTS, AssetHERO, dtTo, dtFrom, 50)

	if err != nil {
		suite.FailNow(err.Error(), "GetTradeHistory")
	}

	suite.NotNil(res)
	//util.Dump("markettrades >", res)
}

func (suite *commonTest) Test_GetLimitOrders() {

	res, err := suite.TestAPI.GetLimitOrders(AssetCNY, AssetBTS, 50)
	if err != nil {
		suite.FailNow(err.Error(), "GetLimitOrders")
	}

	suite.NotNil(res)
	//util.Dump("limitorders >", res)
}

func (suite *commonTest) Test_GetCallOrders() {
	res, err := suite.TestAPI.GetCallOrders(AssetUSD, 50)
	if err != nil {
		suite.FailNow(err.Error(), "GetCallOrders")
	}

	suite.NotNil(res)
	//	util.Dump("callorders >", res)
}

func (suite *commonTest) Test_GetMarginPositions() {
	res, err := suite.TestAPI.GetMarginPositions(UserID2)
	if err != nil {
		suite.FailNow(err.Error(), "GetMarginPositions")
	}

	suite.NotNil(res)
	//util.Dump("marginpositions >", res)
}

func (suite *commonTest) Test_GetSettleOrders() {
	res, err := suite.TestAPI.GetSettleOrders(AssetCNY, 50)
	if err != nil {
		suite.FailNow(err.Error(), "GetSettleOrders")
	}

	suite.NotNil(res)
	//util.Dump("settleorders >", res)
}

func (suite *commonTest) Test_ListAssets() {
	res, err := suite.TestAPI.ListAssets("OPEN.DASH", 2)
	if err != nil {
		suite.FailNow(err.Error(), "ListAssets")
	}

	suite.NotNil(res)
	suite.Len(res, 2)
	//util.Dump("assets >", res)
}

func (suite *commonTest) Test_GetAccountHistory() {

	user := types.NewGrapheneID("1.2.96393")
	start := types.NewGrapheneID("1.11.187698971")
	stop := types.NewGrapheneID("1.11.187658388")

	res, err := suite.TestAPI.GetAccountHistory(user, stop, 30, start)
	if err != nil {
		suite.FailNow(err.Error(), "GetAccountHistory")
	}

	suite.NotNil(res)
	//util.Dump("history >", res)
}

func TestCommon(t *testing.T) {
	testSuite := new(commonTest)
	suite.Run(t, testSuite)
}
