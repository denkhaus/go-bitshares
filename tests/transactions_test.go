package tests

import (
	"testing"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/config"
	"github.com/denkhaus/bitshares/crypto"
	"github.com/denkhaus/bitshares/operations"
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/stretchr/testify/suite"
)

type transactionsAPITest struct {
	suite.Suite
	TestAPI api.BitsharesAPI
	KeyBag  *crypto.KeyBag
	RefTx   *types.Transaction
}

func (suite *transactionsAPITest) SetupTest() {
	suite.TestAPI = NewTestAPI(suite.T(), WsTestApiUrl, RpcTestApiUrl)
	suite.KeyBag = crypto.NewKeyBag()

	if err := suite.KeyBag.Add(TestAccount1PrivKeyOwner); err != nil {
		suite.FailNow(err.Error(), "KeyBag.Add 1")
	}
	if err := suite.KeyBag.Add(TestAccount1PrivKeyActive); err != nil {
		suite.FailNow(err.Error(), "KeyBag.Add 2")
	}

	suite.RefTx = CreateRefTransaction(suite.T())

	if err := suite.TestAPI.WalletUnlock("123456"); err != nil {
		suite.FailNow(err.Error(), "WalletUnlock")
	}
}

func (suite *transactionsAPITest) TearDownTest() {
	if err := suite.TestAPI.WalletLock(); err != nil {
		suite.FailNow(err.Error(), "WalletLock")
	}

	if err := suite.TestAPI.Close(); err != nil {
		suite.FailNow(err.Error(), "Close")
	}
}

func (suite *transactionsAPITest) Test_ChainConfig() {
	res, err := suite.TestAPI.GetChainID()
	if err != nil {
		suite.FailNow(err.Error(), "GetChainID")
	}

	suite.Equal(config.ChainIDTest, res)
}

func (suite *transactionsAPITest) Test_BuildSignedTransaction() {
	op := operations.TransferOperation{
		Extensions: types.Extensions{},
		Amount: types.AssetAmount{
			Amount: 1000,
			Asset:  *AssetTEST,
		},
		From: *TestAccount1ID,
		To:   *TestAccount2ID,
	}

	trx, err := suite.TestAPI.BuildSignedTransaction(suite.KeyBag, AssetTEST, &op)
	if err != nil {
		suite.FailNow(err.Error(), "BuildSignedTransaction")
	}

	//util.Dump("signed trx <", trx)
	suite.compareTransaction(trx, false)
}

// func (suite *transactionsAPITest) Test_SignTransactionCompare() {

// 	op := operations.TransferOperation{
// 		Extensions: types.Extensions{},
// 		Amount: types.AssetAmount{
// 			Amount: 1000,
// 			Asset:  *AssetTEST,
// 		},
// 		From: *TestAccount1ID,
// 		To:   *TestAccount2ID,
// 		Fee: types.AssetAmount{
// 			Amount: 100,
// 			Asset:  *AssetTEST,
// 		},
// 	}

// 	suite.RefTx.Operations = types.Operations{
// 		types.Operation(&op),
// 	}

// 	suite.TestAPI.SetDebug(false)
// 	trxWallet, err := suite.TestAPI.WalletSignTransaction(suite.RefTx, false)
// 	if err != nil {
// 		suite.FailNow(err.Error(), "WalletSignTransaction")
// 	}

// 	suite.compareTransaction(trxWallet, false)
// 	util.Dump("wallet signed trx <", trxWallet)

// 	//sigWallet := trxWallet.Signatures
// 	trxWallet.Signatures = types.Signatures{}
// 	suite.compareTransaction(trxWallet, false)

// 	trxWs, err := suite.TestAPI.SignTransaction(suite.KeyBag, trxWallet)
// 	if err != nil {
// 		suite.FailNow(err.Error(), "SignTransaction")
// 	}

// 	suite.compareTransaction(trxWs, false)
// 	util.Dump("websocket signed trx <", trxWs)

// 	//suite.Equal(sigWallet[0], trxWs.Signatures[0])
// }

func (suite *transactionsAPITest) Test_SignAndVerify() {
	op := operations.TransferOperation{
		Extensions: types.Extensions{},
		Amount: types.AssetAmount{
			Amount: 1000,
			Asset:  *AssetTEST,
		},
		From: *TestAccount1ID,
		To:   *TestAccount2ID,
	}

	trx, err := suite.TestAPI.BuildSignedTransaction(suite.KeyBag, AssetTEST, &op)
	if err != nil {
		suite.FailNow(err.Error(), "BuildSignedTransaction")
	}

	suite.compareTransaction(trx, false)

	v, err := suite.TestAPI.VerifySignedTransaction(suite.KeyBag, trx)
	if err != nil {
		suite.FailNow(err.Error(), "VerifySignedTransaction")
	}

	suite.True(v, "Verified")
}

func (suite *transactionsAPITest) Test_Transfer() {
	op := operations.TransferOperation{
		Extensions: types.Extensions{},
		Amount: types.AssetAmount{
			Amount: 1000,
			Asset:  *AssetTEST,
		},
		From: *TestAccount1ID,
		To:   *TestAccount2ID,
	}

	trx, err := suite.TestAPI.BuildSignedTransaction(suite.KeyBag, AssetTEST, &op)
	if err != nil {
		suite.FailNow(err.Error(), "BuildSignedTransaction")
	}

	suite.compareTransaction(trx, false)

	//suite.TestAPI.SetDebug(true)
	res, err := suite.TestAPI.BroadcastTransaction(trx)
	if err != nil {
		suite.FailNow(err.Error(), "BroadcastTransaction")
	}

	util.Dump("transfer <", res)
}

func (suite *transactionsAPITest) Test_GetAccountBalances() {
	res, err := suite.TestAPI.GetAccountBalances(TestAccount1ID, AssetTEST)
	if err != nil {
		suite.FailNow(err.Error(), "GetAccountBalances 1")
	}

	suite.NotNil(res)
	suite.Len(res, 1)

	//util.Dump("test amount TestAccount1 >", res)

	res, err = suite.TestAPI.GetAccountBalances(TestAccount2ID, AssetTEST)
	if err != nil {
		suite.FailNow(err.Error(), "GetAccountBalances 2")
	}

	suite.NotNil(res)
	suite.Len(res, 1)

	//util.Dump("test amount TestAccount2 >", res)
}

func (suite *transactionsAPITest) compareTransaction(tx *types.Transaction, debug bool) {
	ref, test, err := CompareTransactions(suite.TestAPI, tx, debug)
	if err != nil {
		suite.FailNow(err.Error(), "CompareTransactions")
	}

	suite.Equal(ref, test)
}

func TestTransactionsApi(t *testing.T) {
	testSuite := new(transactionsAPITest)
	suite.Run(t, testSuite)
}
