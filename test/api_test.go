package test

import (
	"log"
	"testing"

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

	asset := objects.NewAsset("1.3.0")
	user := objects.NewAccount("1.2.20067")

	res, err := suite.testAPI.GetAccountBalances(user, asset)
	if err != nil {
		suite.T().Error(err)
	}

	assert.NotNil(suite.T(), res)
	//spew.Dump(res)
}

func (suite *BitsharesAPITest) Test_GetAccounts() {

	accountID := objects.NewGrapheneID("1.2.0")
	res, err := suite.testAPI.GetAccounts(accountID)
	if err != nil {
		suite.T().Error(err)
	}

	assert.NotNil(suite.T(), res)
	assert.Len(suite.T(), res, 1)
	//spew.Dump(res)
}

func (suite *BitsharesAPITest) Test_ListAssets() {
	res, err := suite.testAPI.ListAssets("HERO", 1)
	if err != nil {
		suite.T().Error(err)
	}

	assert.NotNil(suite.T(), res)
	assert.Len(suite.T(), res, 1)
	//spew.Dump(res)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBitsharesAPI(t *testing.T) {
	testSuite := new(BitsharesAPITest)
	suite.Run(t, testSuite)
	testSuite.TearDown()
}
