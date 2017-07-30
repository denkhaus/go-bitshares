package test

import (
	"log"
	"testing"

	"github.com/davecgh/go-spew/spew"
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

func (suite *BitsharesAPITest) Test_GetAccount_Balances() {

	// asset1 := bsapi.NewAsset("1.3.0")
	// user := bsapi.UserAccount{}
	// user.ID = "1.2.20067"

	// res := suite.testAPI.GetAccountBalances(user, asset1)
	// <-res.Done
	// spew.Dump(res)
	//assert.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)
}

func (suite *BitsharesAPITest) Test_GetAccounts() {
	accountID := objects.NewGrapheneID("1.2.0")

	res, err := suite.testAPI.GetAccounts(accountID)
	if err != nil {
		suite.Error(err)
	}

	assert.NotNil(suite.T(), res)
	spew.Dump(res)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBitsharesAPI(t *testing.T) {
	suite.Run(t, new(BitsharesAPITest))
}
