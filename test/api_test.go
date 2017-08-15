package test

import (
	"testing"
	"time"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/objects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	testURL = "wss://bitshares.openledger.info/ws"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type bitsharesAPITest struct {
	suite.Suite
	TestAPI          api.BitsharesAPI
	AssetCNY         *objects.GrapheneID
	AssetBTS         *objects.GrapheneID
	AssetUSD         *objects.GrapheneID
	BitAssetDataCNY  *objects.GrapheneID
	UserID1          *objects.GrapheneID
	UserID2          *objects.GrapheneID
	UserID3          *objects.GrapheneID
	LimitOrder1      *objects.GrapheneID
	CallOrder1       *objects.GrapheneID
	SettleOrder1     *objects.GrapheneID
	ChainIDBitShares string
}

func (suite *bitsharesAPITest) SetupTest() {
	suite.ChainIDBitShares = "4018d7844c78f6a6c41c6a552b898022310fc5dec06da467ee7905a8dad512c8"
	suite.UserID1 = objects.NewGrapheneID("1.2.282")          //xeroc user account
	suite.UserID2 = objects.NewGrapheneID("1.2.253")          //stan user account
	suite.UserID3 = objects.NewGrapheneID("1.2.0")            //committee-account user account
	suite.AssetCNY = objects.NewGrapheneID("1.3.113")         //cny asset
	suite.AssetBTS = objects.NewGrapheneID("1.3.0")           //bts asset
	suite.AssetUSD = objects.NewGrapheneID("1.3.121")         // usd asset
	suite.BitAssetDataCNY = objects.NewGrapheneID("2.4.13")   //cny bitasset data id
	suite.LimitOrder1 = objects.NewGrapheneID("1.7.22765740") // random LimitOrder ObjectID
	suite.CallOrder1 = objects.NewGrapheneID("1.8.4582")      // random CallOrder ObjectID
	suite.SettleOrder1 = objects.NewGrapheneID("1.4.1655")    // random SettleOrder ObjectID

	api := api.New(testURL)
	if err := api.Connect(); err != nil {
		suite.Fail(err.Error(), "Connect")
	}

	suite.TestAPI = api
}

func (suite *bitsharesAPITest) TearDown() {
	suite.TestAPI.Close()
}

func (suite *bitsharesAPITest) Test_GetChainID() {

	res, err := suite.TestAPI.GetChainID()
	if err != nil {
		suite.Fail(err.Error(), "GetChainID")
	}

	assert.Equal(suite.T(), res, suite.ChainIDBitShares)
}

func (suite *bitsharesAPITest) Test_GetAccountBalances() {

	res, err := suite.TestAPI.GetAccountBalances(suite.UserID2, suite.AssetBTS)
	if err != nil {
		suite.Fail(err.Error(), "GetAccountBalances")
	}

	assert.NotNil(suite.T(), res)
	//spew.Dump(res)
}

func (suite *bitsharesAPITest) Test_GetAccounts() {

	res, err := suite.TestAPI.GetAccounts(suite.UserID3)
	if err != nil {
		suite.Fail(err.Error(), "GetAccounts")
	}

	assert.NotNil(suite.T(), res)
	assert.Len(suite.T(), res, 1)
	//spew.Dump(res)
}

func (suite *bitsharesAPITest) Test_GetObjects() {

	res, err := suite.TestAPI.GetObjects(
		suite.UserID1,
		suite.AssetCNY,
		suite.BitAssetDataCNY,
		suite.LimitOrder1,
		suite.CallOrder1,
		suite.SettleOrder1,
	)

	if err != nil {
		suite.Fail(err.Error(), "GetObjects")
	}

	assert.NotNil(suite.T(), res)
	assert.Len(suite.T(), res, 6)
	//util.Dump("objects out", res)
}

func (suite *bitsharesAPITest) Test_GetAccountByName() {

	res, err := suite.TestAPI.GetAccountByName("openledger")
	if err != nil {
		suite.Fail(err.Error(), "GetAccountByName")
	}

	assert.NotNil(suite.T(), res)
	//util.Dump("accounts out", res)
}

func (suite *bitsharesAPITest) Test_GetTradeHistory() {
	dtTo := time.Now().UTC()
	dtFrom := dtTo.Add(-time.Hour)

	res, err := suite.TestAPI.GetTradeHistory("CNY", "BTS", dtTo, dtFrom, 50)
	if err != nil {
		suite.Fail(err.Error(), "GetTradeHistory")
	}

	assert.NotNil(suite.T(), res)
	//util.Dump("markettrades out", res)
}

func (suite *bitsharesAPITest) Test_GetLimitOrders() {

	res, err := suite.TestAPI.GetLimitOrders(suite.AssetCNY, suite.AssetBTS, 50)
	if err != nil {
		suite.Fail(err.Error(), "GetLimitOrders")
	}

	assert.NotNil(suite.T(), res)
	//util.Dump("limitorders out", res)
}

func (suite *bitsharesAPITest) Test_GetCallOrders() {

	res, err := suite.TestAPI.GetCallOrders(suite.AssetUSD, 50)
	if err != nil {
		suite.Fail(err.Error(), "GetCallOrders")
	}

	assert.NotNil(suite.T(), res)
	//	util.Dump("callorders out", res)
}

func (suite *bitsharesAPITest) Test_GetSettleOrders() {

	res, err := suite.TestAPI.GetSettleOrders(suite.AssetCNY, 50)
	if err != nil {
		suite.Fail(err.Error(), "GetSettleOrders")
	}

	assert.NotNil(suite.T(), res)
	//util.Dump("settleorders out", res)
}

func (suite *bitsharesAPITest) Test_ListAssets() {
	res, err := suite.TestAPI.ListAssets("HERO", 2)
	if err != nil {
		suite.Fail(err.Error(), "ListAssets")
	}

	assert.NotNil(suite.T(), res)
	assert.Len(suite.T(), res, 2)
	//util.Dump("assets out", res)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBitsharesAPI(t *testing.T) {
	testSuite := new(bitsharesAPITest)
	suite.Run(t, testSuite)
	testSuite.TearDown()
}
