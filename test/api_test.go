package test

import (
	"log"
	"testing"
	"time"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/objects"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	testURL = "wss://bitshares.openledger.info/ws"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type BitsharesAPITest struct {
	suite.Suite
	testAPI *api.BitsharesApi
}

func (suite *BitsharesAPITest) SetupTest() {
	testAPI, err := api.New(testURL)
	if err != nil {
		log.Fatal(err)
	}

	suite.testAPI = testAPI
}

func (suite *BitsharesAPITest) TearDown() {
	suite.testAPI.Close()
}

func (suite *BitsharesAPITest) Test_GetAccountBalances() {

	asset := objects.NewGrapheneID("1.3.0")
	user := objects.NewGrapheneID("1.2.20067")

	res, err := suite.testAPI.GetAccountBalances(user, asset)
	if err != nil {
		suite.T().Error(errors.Annotate(err, "GetAccountBalances"))
	}

	assert.NotNil(suite.T(), res)
	//spew.Dump(res)
}

func (suite *BitsharesAPITest) Test_GetAccounts() {

	accountID := objects.NewGrapheneID("1.2.0")
	res, err := suite.testAPI.GetAccounts(accountID)
	if err != nil {
		suite.T().Error(errors.Annotate(err, "GetAccounts"))
	}

	assert.NotNil(suite.T(), res)
	assert.Len(suite.T(), res, 1)
	//spew.Dump(res)
}

func (suite *BitsharesAPITest) Test_GetObjects() {

	ID1 := objects.NewGrapheneID("1.2.282") //xeroc user account
	ID2 := objects.NewGrapheneID("1.3.113") //cny asset
	ID3 := objects.NewGrapheneID("2.4.13")  //cny bitasset data id

	res, err := suite.testAPI.GetObjects(ID1, ID2, ID3)
	if err != nil {
		suite.T().Error(errors.Annotate(err, "GetObjects"))
	}

	assert.NotNil(suite.T(), res)
	assert.Len(suite.T(), res, 3)
	//util.Dump("objects out", res)
}

func (suite *BitsharesAPITest) Test_GetAccountByName() {

	res, err := suite.testAPI.GetAccountByName("openledger")
	if err != nil {
		suite.T().Error(errors.Annotate(err, "GetAccountByName"))
	}

	assert.NotNil(suite.T(), res)
	//util.Dump("accounts out", res)
}

func (suite *BitsharesAPITest) Test_GetTradeHistory() {
	dtTo := time.Now().UTC()
	dtFrom := dtTo.Add(-time.Hour)

	res, err := suite.testAPI.GetTradeHistory("CNY", "BTS", dtTo, dtFrom, 50)
	if err != nil {
		suite.T().Error(errors.Annotate(err, "GetTradeHistory"))
	}

	assert.NotNil(suite.T(), res)
	//util.Dump("markettrades out", res)
}

func (suite *BitsharesAPITest) Test_GetLimitOrders() {

	cny := objects.NewGrapheneID("1.3.113") //cny asset
	bts := objects.NewGrapheneID("1.3.0")   //bts asset

	res, err := suite.testAPI.GetLimitOrders(cny, bts, 50)
	if err != nil {
		suite.T().Error(errors.Annotate(err, "GetLimitOrders"))
	}

	assert.NotNil(suite.T(), res)
	util.Dump("limitorders out", res)
}

func (suite *BitsharesAPITest) Test_ListAssets() {
	res, err := suite.testAPI.ListAssets("HERO", 2)
	if err != nil {
		suite.T().Error(errors.Annotate(err, "ListAssets"))
	}

	assert.NotNil(suite.T(), res)
	assert.Len(suite.T(), res, 2)
	//util.Dump("assets out", res)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBitsharesAPI(t *testing.T) {
	testSuite := new(BitsharesAPITest)
	suite.Run(t, testSuite)
	testSuite.TearDown()
}
