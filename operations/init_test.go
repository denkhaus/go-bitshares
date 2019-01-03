package operations

import (
	"fmt"
	"testing"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/crypto"
	"github.com/denkhaus/bitshares/tests"
	"github.com/denkhaus/bitshares/types"
	"github.com/stretchr/testify/suite"

	// importing this initializes sample data fetching
	_ "github.com/denkhaus/bitshares/gen/samples"
)

type operationsAPITest struct {
	suite.Suite
	TestAPI api.BitsharesAPI
	RefTx   *types.SignedTransaction
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
	suite.compareTransaction(0, suite.RefTx, false)
}

func (suite *operationsAPITest) Test_WalletSerializeTransaction() {
	hex, err := suite.TestAPI.WalletSerializeTransaction(suite.RefTx)
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
			OperationFee: types.OperationFee{
				Fee: &types.AssetAmount{
					Amount: 100,
					Asset:  *types.NewGrapheneID("1.3.0"),
				},
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
			Extensions:     types.CallOrderUpdateExtensions{},
		},
	}

	if err := suite.TestAPI.SignWithKeys(keyBag.Privates(), suite.RefTx); err != nil {
		suite.FailNow(err.Error(), "SignTransaction")
	}

	suite.compareTransaction(0, suite.RefTx, false)
}

func (suite *operationsAPITest) compareTransaction(sampleIdx int, tx *types.SignedTransaction, debug bool) {
	ref, test, err := tests.CompareTransactions(suite.TestAPI, tx, debug)
	if err != nil {
		suite.FailNow(err.Error(), "compareTransactions")
	}

	suite.Equal(
		ref,
		test,
		fmt.Sprintf("on sample index %d\n", sampleIdx),
	)
}

func TestOperations(t *testing.T) {
	testSuite := new(operationsAPITest)
	suite.Run(t, testSuite)
}
