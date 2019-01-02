package operations

import (
	"testing"

	"github.com/denkhaus/bitshares"
	"github.com/denkhaus/bitshares/crypto"
	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/tests"
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/logging"
	"github.com/stretchr/testify/suite"

	// importing this initializes sample data fetching
	_ "github.com/denkhaus/bitshares/gen/samples"
)

type operationsAPITest struct {
	suite.Suite
	WebsocketAPI bitshares.WebsocketAPI
	WalletAPI    bitshares.WalletAPI
	RefTx        *types.SignedTransaction
}

func (suite *operationsAPITest) SetupTest() {
	suite.WebsocketAPI = tests.NewWebsocketTestAPI(
		suite.T(),
		tests.WsFullApiUrl,
	)
	suite.WalletAPI = tests.NewWalletTestAPI(
		suite.T(),
		tests.RpcFullApiUrl,
	)
	suite.RefTx = tests.CreateRefTransaction(suite.T())
}

func (suite *operationsAPITest) TearDownTest() {
	if err := suite.WebsocketAPI.Close(); err != nil {
		suite.FailNow(err.Error(), "Close")
	}
}

func (suite *operationsAPITest) OpSamplesTest(op types.Operation) {
	samples, err := data.GetSamplesByType(op.Type())
	if err != nil {
		if err == data.ErrNoSampleDataAvailable {
			logging.Warnf("no sample data available for %s", op.Type())
			return
		}
		suite.FailNow(err.Error(), "GetSamplesByType")
	}

	um, ok := op.(types.Unmarshalable)
	if !ok {
		suite.FailNow("test error", "operation %v is not unmarshalable", op)
	}

	for idx, sample := range samples {
		if err := um.UnmarshalJSON([]byte(sample)); err != nil {
			suite.FailNow(err.Error(), "UnmarshalJSON")
		}

		suite.RefTx.Operations = types.Operations{
			op,
		}

		suite.compareTransaction(idx, suite.RefTx, false)
	}
}

func (suite *operationsAPITest) Test_SerializeRefTransaction() {
	suite.compareTransaction(0, suite.RefTx, false)
}

func (suite *operationsAPITest) Test_WalletSerializeTransaction() {
	hex, err := suite.WalletAPI.SerializeTransaction(suite.RefTx)
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
					Asset:  types.AssetIDFromObject(tests.AssetBTS),
				},
			},
			DeltaDebt: types.AssetAmount{
				Amount: 10000,
				Asset:  types.AssetIDFromObject(tests.AssetUSD),
			},
			DeltaCollateral: types.AssetAmount{
				Amount: 100000000,
				Asset:  types.AssetIDFromObject(tests.AssetBTS),
			},

			FundingAccount: types.AccountIDFromObject(tests.UserID3),
			Extensions:     types.CallOrderUpdateExtensions{},
		},
	}

	if err := crypto.SignWithKeys(keyBag.Privates(), suite.RefTx); err != nil {
		suite.FailNow(err.Error(), "SignWithKeys")
	}

	suite.compareTransaction(0, suite.RefTx, false)
}

func (suite *operationsAPITest) compareTransaction(sampleIdx int, tx *types.SignedTransaction, debug bool) {
	ref, test, err := tests.CompareTransactions(suite.WalletAPI, tx, debug)
	if err != nil {
		suite.FailNow(err.Error(), "compareTransaction")
	}

	suite.Equal(
		ref,
		test,
		"on sample index %d",
		sampleIdx,
	)
}

func TestOperations(t *testing.T) {
	testSuite := new(operationsAPITest)
	suite.Run(t, testSuite)
}
