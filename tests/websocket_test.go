package tests

import (
	"testing"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/config"
	"github.com/denkhaus/bitshares/crypto"
	"github.com/denkhaus/bitshares/operations"
	"github.com/denkhaus/bitshares/types"
	"github.com/stretchr/testify/suite"
)

type websocketAPITest struct {
	suite.Suite
	WebsocketAPI api.WebsocketAPI
	WalletAPI    api.WalletAPI
	KeyBag       *crypto.KeyBag
}

func (suite *websocketAPITest) SetupSuite() {
	suite.WebsocketAPI = NewWebsocketTestAPI(suite.T(), WsTestApiUrl)
	suite.WalletAPI = NewWalletTestAPI(suite.T(), RpcTestApiUrl, config.ChainIDTest)

	suite.KeyBag = crypto.NewKeyBag()
	if err := suite.KeyBag.Add(TestAccount1PrivKeyActive); err != nil {
		suite.FailNow(err.Error(), "KeyBag.Add 1")
	}

	if err := suite.KeyBag.Add(TestAccount2PrivKeyActive); err != nil {
		suite.FailNow(err.Error(), "KeyBag.Add 2")
	}
}

func (suite *websocketAPITest) TearDownSuite() {
	if err := suite.WebsocketAPI.Close(); err != nil {
		suite.FailNow(err.Error(), "Close [ws]")
	}
	if err := suite.WalletAPI.Close(); err != nil {
		suite.FailNow(err.Error(), "Close [wallet]")
	}
}

func (suite *websocketAPITest) Test_ChainConfig() {
	res, err := suite.WebsocketAPI.GetChainID()
	if err != nil {
		suite.FailNow(err.Error(), "GetChainID")
	}

	suite.Equal(config.ChainIDTest, res)
}

func (suite *websocketAPITest) Test_BuildSignedTransaction() {
	op := operations.TransferOperation{
		Extensions: types.Extensions{},
		Amount: types.AssetAmount{
			Amount: 1000,
			Asset:  *AssetTEST,
		},
		From: *TestAccount2ID,
		To:   *TestAccount1ID,
	}

	trx, err := suite.WebsocketAPI.BuildSignedTransaction(suite.KeyBag, AssetTEST, &op)
	if err != nil {
		suite.FailNow(err.Error(), "BuildSignedTransaction")
	}

	//util.Dump("signed trx <", trx)
	suite.compareTransaction(trx, false)
}

func (suite *websocketAPITest) Test_TransferExtended() {
	props, err := suite.WebsocketAPI.GetDynamicGlobalProperties()
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

	fees, err := suite.WebsocketAPI.GetRequiredFees(trx.Operations, AssetTEST)
	if err != nil {
		suite.FailNow(err.Error(), "GetRequiredFees")
	}

	if err := trx.Operations.ApplyFees(fees); err != nil {
		suite.FailNow(err.Error(), "ApplyFees")
	}

	suite.compareTransaction(trx, false)

	res, err := suite.WalletAPI.SignTransaction(trx, true)
	if err != nil {
		suite.FailNow(err.Error(), "SignTransaction")
	}

	_ = res
	//logging.Dump("transfer <", res)
}

func (suite *websocketAPITest) Test_SignAndVerify() {
	op := operations.TransferOperation{
		Extensions: types.Extensions{},
		Amount: types.AssetAmount{
			Amount: 1000,
			Asset:  *AssetTEST,
		},
		From: *TestAccount2ID,
		To:   *TestAccount1ID,
	}

	trx, err := suite.WebsocketAPI.BuildSignedTransaction(suite.KeyBag, AssetTEST, &op)
	if err != nil {
		suite.FailNow(err.Error(), "BuildSignedTransaction")
	}

	suite.compareTransaction(trx, false)

	v, err := crypto.VerifySignedTransaction(suite.KeyBag, trx)
	if err != nil {
		suite.FailNow(err.Error(), "VerifySignedTransaction")
	}

	suite.True(v, "Verified")
}

func (suite *websocketAPITest) Test_Transfer() {
	op := operations.TransferOperation{
		Extensions: types.Extensions{},
		Amount: types.AssetAmount{
			Amount: 100000,
			Asset:  *AssetTEST,
		},
		From: *TestAccount2ID,
		To:   *TestAccount1ID,
	}

	trx, err := suite.WebsocketAPI.BuildSignedTransaction(suite.KeyBag, AssetTEST, &op)
	if err != nil {
		suite.FailNow(err.Error(), "BuildSignedTransaction")
	}

	//util.Dump("trx websocket <", trx)
	suite.compareTransaction(trx, false)

	//suite.TestAPI.SetDebug(true)
	if err := suite.WebsocketAPI.BroadcastTransaction(trx); err != nil {
		suite.FailNow(err.Error(), "BroadcastTransaction")
	}
}

func (suite *websocketAPITest) Test_GetAccountBalances() {
	res, err := suite.WebsocketAPI.GetAccountBalances(TestAccount1ID, AssetTEST)
	if err != nil {
		suite.FailNow(err.Error(), "GetAccountBalances 1")
	}

	suite.NotNil(res)
	suite.Len(res, 1)

	//util.Dump("test amount TestAccount1 >", res)

	res, err = suite.WebsocketAPI.GetAccountBalances(TestAccount2ID, AssetTEST)
	if err != nil {
		suite.FailNow(err.Error(), "GetAccountBalances 2")
	}

	suite.NotNil(res)
	suite.Len(res, 1)

	//util.Dump("test amount TestAccount2 >", res)
}

func (suite *websocketAPITest) compareTransaction(tx *types.SignedTransaction, debug bool) {
	ref, test, err := CompareTransactions(suite.WalletAPI, tx, debug)
	if err != nil {
		suite.FailNow(err.Error(), "Compare Transactions")
	}

	suite.Equal(ref, test)
}

func TestWebsocketAPI(t *testing.T) {
	testSuite := new(websocketAPITest)
	suite.Run(t, testSuite)
}
