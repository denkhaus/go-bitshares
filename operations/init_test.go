package operations

import (
	"testing"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/tests"
	"github.com/denkhaus/bitshares/types"

	"github.com/stretchr/testify/suite"

	// importing this initializes sample data fetching
	_ "github.com/denkhaus/bitshares/gen/samples"
)

//Note: operation tests may fail for now cause extensions marshalling is questionable.
type operationsAPITest struct {
	suite.Suite
	TestAPI api.BitsharesAPI
	RefTx   *types.Transaction
}

func (suite *operationsAPITest) SetupTest() {
	api := api.New(tests.WsFullApiUrl, tests.RpcApiUrl)

	if err := api.Connect(); err != nil {
		suite.FailNow(err.Error(), "Connect")
	}

	api.OnError(func(err error) {
		suite.FailNow(err.Error(), "OnError")
	})

	suite.TestAPI = api

	tx := types.NewTransaction()
	tx.RefBlockNum = 34294
	tx.RefBlockPrefix = 3707022213

	if err := tx.Expiration.UnmarshalJSON([]byte(`"2016-04-06T08:29:27"`)); err != nil {
		suite.FailNow(err.Error(), "Unmarshal expiration")
	}

	suite.RefTx = tx
}

func (suite *operationsAPITest) TearDownTest() {
	if err := suite.TestAPI.Close(); err != nil {
		suite.FailNow(err.Error(), "Close")
	}
}

func (suite *operationsAPITest) Test_SerializeEmptyTransaction() {

	tx := types.NewTransaction()
	if err := tx.Expiration.UnmarshalJSON([]byte(`"2016-04-06T08:29:27"`)); err != nil {
		suite.FailNow(err.Error(), "Unmarshal expiration")
	}

	suite.compareTransaction(tx, false)
}

func (suite *operationsAPITest) Test_SerializeTransaction() {
	hex, err := suite.TestAPI.SerializeTransaction(suite.RefTx)
	if err != nil {
		suite.FailNow(err.Error(), "SerializeTransaction")
	}

	suite.NotNil(hex)
	suite.Equal("f68585abf4dce7c80457000000", hex)
}

func (suite *operationsAPITest) compareTransaction(tx *types.Transaction, debug bool) {
	ref, test, err := tests.CompareTransactions(suite.TestAPI, tx, debug)
	if err != nil {
		suite.FailNow(err.Error(), "compareTry")
	}

	suite.Equal(ref, test)
}

func TestOperations(t *testing.T) {
	testSuite := new(operationsAPITest)
	suite.Run(t, testSuite)
}
