package tests

import (
	"testing"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/crypto"
	"github.com/denkhaus/bitshares/operations"
	"github.com/denkhaus/bitshares/types"
	"github.com/stretchr/testify/suite"
)

type walletAPITest struct {
	suite.Suite
	TestAPI api.BitsharesAPI
	KeyBag  *crypto.KeyBag
}

func (suite *walletAPITest) SetupSuite() {
	suite.TestAPI = NewTestAPI(suite.T(), WsTestApiUrl, RpcTestApiUrl)

	suite.KeyBag = crypto.NewKeyBag()
	if err := suite.KeyBag.Add(TestAccount1PrivKeyActive); err != nil {
		suite.FailNow(err.Error(), "KeyBag.Add 1")
	}

	if err := suite.KeyBag.Add(TestAccount2PrivKeyActive); err != nil {
		suite.FailNow(err.Error(), "KeyBag.Add 2")
	}

	if err := suite.TestAPI.WalletUnlock("123456"); err != nil {
		suite.FailNow(err.Error(), "WalletUnlock")
	}
}

func (suite *walletAPITest) TearDownSuite() {
	if err := suite.TestAPI.WalletLock(); err != nil {
		suite.FailNow(err.Error(), "WalletLock")
	}

	if err := suite.TestAPI.Close(); err != nil {
		suite.FailNow(err.Error(), "Close")
	}
}

func (suite *walletAPITest) Test_TransferExtended() {
	props, err := suite.TestAPI.GetDynamicGlobalProperties()
	if err != nil {
		suite.FailNow(err.Error(), "GetDynamicGlobalProperties")
	}

	trx, err := types.NewSignedTransactionWithBlockData(props)
	if err != nil {
		suite.FailNow(err.Error(), "NewSignedTransactionWithBlockData")
	}

	trx.Operations = types.Operations{
		&operations.TransferOperation{
			Extensions: types.Extensions{},
			Amount: types.AssetAmount{
				Amount: 100000,
				Asset:  *AssetTEST,
			},
			From: *TestAccount1ID,
			To:   *TestAccount2ID,
		},
	}

	// logging.SetDebug(true)
	// defer logging.SetDebug(false)

	fees, err := suite.TestAPI.GetRequiredFees(trx.Operations, AssetTEST)
	if err != nil {
		suite.FailNow(err.Error(), "GetRequiredFees")
	}

	if err := trx.Operations.ApplyFees(fees); err != nil {
		suite.FailNow(err.Error(), "ApplyFees")
	}

	suite.compareTransaction(trx, false)

	res, err := suite.TestAPI.WalletSignTransaction(trx, true)
	if err != nil {
		suite.FailNow(err.Error(), "WalletSignTransaction")
	}

	_ = res
	//util.Dump("transfer <", res)
}

func (suite *walletAPITest) compareTransaction(tx *types.SignedTransaction, debug bool) {
	ref, test, err := CompareTransactions(suite.TestAPI, tx, debug)
	if err != nil {
		suite.FailNow(err.Error(), "Compare Transactions")
	}

	suite.Equal(ref, test)
}

func TestWalletAPI(t *testing.T) {
	testSuite := new(walletAPITest)
	suite.Run(t, testSuite)
}
