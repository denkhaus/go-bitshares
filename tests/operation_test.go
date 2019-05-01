package tests

import (
	"testing"

	"github.com/denkhaus/bitshares"
	"github.com/denkhaus/bitshares/crypto"
	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/operations"
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
	suite.WebsocketAPI = NewWebsocketTestAPI(
		suite.T(),
		WsFullApiUrl,
	)
	suite.WalletAPI = NewWalletTestAPI(
		suite.T(),
		RpcFullApiUrl,
	)
	suite.RefTx = CreateRefTransaction(suite.T())
}

func (suite *operationsAPITest) TearDownTest() {
	if err := suite.WebsocketAPI.Close(); err != nil {
		suite.FailNow(err.Error(), "Close")
	}
}
func (suite *operationsAPITest) Test_AccountUpdateOperation() {
	suite.samplesTest(&operations.AccountUpdateOperation{})
}

func (suite *operationsAPITest) Test_AssetPublishFeedOperation() {
	op := operations.AssetPublishFeedOperation{
		Extensions: types.Extensions{},
	}

	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_AssetFundFeePoolOperation() {
	op := operations.AssetFundFeePoolOperation{
		Extensions: types.Extensions{},
	}

	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_AssetUpdateFeedProducersOperation() {
	op := operations.AssetUpdateFeedProducersOperation{
		Extensions: types.Extensions{},
	}

	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_CustomOperation() {
	suite.samplesTest(&operations.CustomOperation{})
}

func (suite *operationsAPITest) Test_CallOrderUpdateOperation() {
	suite.samplesTest(&operations.CallOrderUpdateOperation{})
}

func (suite *operationsAPITest) Test_BalanceClaimOperation() {
	suite.samplesTest(&operations.BalanceClaimOperation{})
}

func (suite *operationsAPITest) Test_AssetUpdateBitassetOperation() {
	op := operations.AssetUpdateBitassetOperation{
		Extensions: types.Extensions{},
	}

	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_AssetReserveOperation() {
	op := operations.AssetReserveOperation{
		Extensions: types.Extensions{},
	}

	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_WithdrawPermissionUpdateOperation() {
	suite.samplesTest(&operations.WithdrawPermissionUpdateOperation{})
}

func (suite *operationsAPITest) Test_TransferFromBlindOperation() {
	suite.samplesTest(&operations.TransferFromBlindOperation{})
}

func (suite *operationsAPITest) Test_VestingBalanceWithdrawOperation() {
	suite.samplesTest(&operations.VestingBalanceWithdrawOperation{})
}

func (suite *operationsAPITest) Test_TransferToBlindOperation() {
	suite.samplesTest(&operations.TransferToBlindOperation{})
}

func (suite *operationsAPITest) Test_WitnessCreateOperation() {
	suite.samplesTest(&operations.WitnessCreateOperation{})
}

func (suite *operationsAPITest) Test_TransferOperation() {
	op := operations.TransferOperation{
		Extensions: types.Extensions{},
	}

	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_AssetClaimFeesOperation() {
	suite.samplesTest(&operations.AssetClaimFeesOperation{})
}

func (suite *operationsAPITest) Test_AssetGlobalSettleOperation() {
	suite.samplesTest(&operations.AssetGlobalSettleOperation{})
}

func (suite *operationsAPITest) Test_WithdrawPermissionDeleteOperation() {
	suite.samplesTest(&operations.WithdrawPermissionDeleteOperation{})
}

func (suite *operationsAPITest) Test_ProposalUpdateOperation() {
	suite.samplesTest(&operations.ProposalUpdateOperation{})
}

func (suite *operationsAPITest) Test_WorkerCreateOperation() {
	suite.samplesTest(&operations.WorkerCreateOperation{})
}

func (suite *operationsAPITest) Test_WitnessUpdateOperation() {
	op := operations.WitnessUpdateOperation{}
	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_VestingBalanceCreateOperation() {
	suite.samplesTest(&operations.VestingBalanceCreateOperation{})
}

func (suite *operationsAPITest) Test_AccountWhitelistOperation() {
	op := operations.AccountWhitelistOperation{
		Extensions: types.Extensions{},
	}

	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_ProposalDeleteOperation() {
	op := operations.ProposalDeleteOperation{}

	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_ProposalCreateOperation() {
	op := operations.ProposalCreateOperation{
		Extensions: types.Extensions{},
	}

	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_AssertOperation() {
	suite.samplesTest(&operations.AssertOperation{})
}

func (suite *operationsAPITest) Test_AssetSettleOperation() {
	op := operations.AssetSettleOperation{
		Extensions: types.Extensions{},
	}

	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_AccountUpgradeOperation() {
	op := operations.AccountUpgradeOperation{
		Extensions: types.Extensions{},
	}

	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_WithdrawPermissionCreateOperation() {
	suite.samplesTest(&operations.WithdrawPermissionCreateOperation{})
}

func (suite *operationsAPITest) Test_WithdrawPermissionClaimOperation() {
	suite.samplesTest(&operations.WithdrawPermissionClaimOperation{})
}

func (suite *operationsAPITest) Test_OverrideTransferOperation() {
	op := operations.OverrideTransferOperation{
		Extensions: types.Extensions{},
	}

	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_AssetCreateOperation() {
	op := operations.AssetCreateOperation{
		Extensions: types.Extensions{},
	}

	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_FillOrderOperation() {
	suite.samplesTest(&operations.FillOrderOperation{})
}

func (suite *operationsAPITest) Test_BidCollateralOperation() {
	suite.samplesTest(&operations.BidCollateralOperation{})
}
func (suite *operationsAPITest) Test_LimitOrderCreateOperation() {
	op := operations.LimitOrderCreateOperation{
		Extensions: types.Extensions{},
	}

	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_AssetUpdateOperation() {
	op := operations.AssetUpdateOperation{
		Extensions: types.Extensions{},
	}

	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_AccountTransferOperation() {
	suite.samplesTest(&operations.AccountTransferOperation{})
}

func (suite *operationsAPITest) Test_AssetIssueOperation() {
	op := operations.AssetIssueOperation{
		Extensions: types.Extensions{},
	}

	suite.samplesTest(&op)
}

func (suite *operationsAPITest) Test_CommitteeMemberCreateOperation() {
	suite.samplesTest(&operations.CommitteeMemberCreateOperation{})
}

func (suite *operationsAPITest) Test_CommitteeMemberUpdateOperation() {
	suite.samplesTest(&operations.CommitteeMemberUpdateOperation{})
}

func (suite *operationsAPITest) Test_AccountCreateOperation() {
	suite.samplesTest(&operations.AccountCreateOperation{})
}

func (suite *operationsAPITest) Test_LimitOrderCancelOperation() {
	op := operations.LimitOrderCancelOperation{
		Extensions: types.Extensions{},
	}

	suite.samplesTest(&op)
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
		&operations.CallOrderUpdateOperation{
			OperationFee: types.OperationFee{
				Fee: &types.AssetAmount{
					Amount: 100,
					Asset:  types.AssetIDFromObject(AssetBTS),
				},
			},
			DeltaDebt: types.AssetAmount{
				Amount: 10000,
				Asset:  types.AssetIDFromObject(AssetUSD),
			},
			DeltaCollateral: types.AssetAmount{
				Amount: 100000000,
				Asset:  types.AssetIDFromObject(AssetBTS),
			},

			FundingAccount: types.AccountIDFromObject(UserID3),
			Extensions:     types.CallOrderUpdateExtensions{},
		},
	}

	if err := crypto.SignWithKeys(keyBag.Privates(), suite.RefTx); err != nil {
		suite.FailNow(err.Error(), "SignWithKeys")
	}

	suite.compareTransaction(0, suite.RefTx, false)
}

func (suite *operationsAPITest) samplesTest(op types.Operation) {
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

func (suite *operationsAPITest) compareTransaction(sampleIdx int, tx *types.SignedTransaction, debug bool) {
	ref, test, err := CompareTransactions(suite.WalletAPI, tx, debug)
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
