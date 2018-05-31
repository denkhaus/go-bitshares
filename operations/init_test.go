package operations

import (
	"testing"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/crypto"
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
	suite.TestAPI = tests.NewTestAPI(
		suite.T(),
		tests.WsFullApiUrl,
		tests.RpcFullApiUrl,
	)

	suite.RefTx = tests.CreateRefTransaction(suite.T())
}

func (suite *operationsAPITest) TearDownTest() {
	if err := suite.TestAPI.Close(); err != nil {
		suite.FailNow(err.Error(), "Close")
	}
}

func (suite *operationsAPITest) Test_SerializeRefTransaction() {
	suite.compareTransaction(suite.RefTx, false)
}

func (suite *operationsAPITest) Test_SerializeTransaction() {
	hex, err := suite.TestAPI.SerializeTransaction(suite.RefTx)
	if err != nil {
		suite.FailNow(err.Error(), "SerializeTransaction")
	}

	suite.NotNil(hex)
	suite.Equal("f68585abf4dce7c80457000000", hex)
}

func (suite *operationsAPITest) Test_SampleOperation() {

	TestWIF := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"

	keyBag := crypto.NewKeyBag()
	if err := keyBag.Add(TestWIF); err != nil {
		suite.FailNow(err.Error(), "KeyBag.Add")
	}

	suite.RefTx.Operations = types.Operations{
		&CallOrderUpdateOperation{
			Fee: types.AssetAmount{
				Amount: 100,
				Asset:  *types.NewGrapheneID("1.3.0"),
			},
			DeltaDebt: types.AssetAmount{
				Amount: 10000,
				Asset:  *types.NewGrapheneID("1.3.22"),
			},
			DeltaCollateral: types.AssetAmount{
				Amount: 100000000,
				Asset:  *types.NewGrapheneID("1.3.0"),
			},

			FundingAccount: *types.NewGrapheneID("1.2.29"),
			Extensions:     types.Extensions{},
		},
	}

	trx, err := suite.TestAPI.SignWithKeys(keyBag.Privates(), suite.RefTx)
	if err != nil {
		suite.FailNow(err.Error(), "SignTransaction")
	}

	suite.NotNil(trx)

	expected := "f68585abf4dce7c8045701036400000000000000001d00e1f" +
		"50500000000001027000000000000160000011f2627efb5c5" +
		"144440e06ff567f1a09928d699ac6f5122653cd7173362a1a" +
		"e20205952c874ed14ccec050be1c86c1a300811763ef3b481" +
		"e562e0933c09b40e31fb"

	test := trx.ToHex()

	suite.Equal(expected[:len(expected)-130], test[:len(test)-130])
}

func (suite *operationsAPITest) compareTransaction(tx *types.Transaction, debug bool) {
	ref, test, err := tests.CompareTransactions(suite.TestAPI, tx, debug)
	if err != nil {
		suite.FailNow(err.Error(), "compareTransactions")
	}

	suite.Equal(ref, test)
}

func TestOperations(t *testing.T) {
	testSuite := new(operationsAPITest)
	suite.Run(t, testSuite)
}
