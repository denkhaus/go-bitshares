package operations

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/denkhaus/bitshares/api"
	"github.com/denkhaus/bitshares/tests"
	"github.com/denkhaus/bitshares/types"

	"github.com/denkhaus/bitshares/util"
	"github.com/stretchr/testify/suite"
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
		suite.Fail(err.Error(), "Connect")
	}

	api.OnError(func(err error) {
		suite.Fail(err.Error(), "OnError")
	})

	suite.TestAPI = api

	tx := types.NewTransaction()
	tx.RefBlockNum = 34294
	tx.RefBlockPrefix = 3707022213

	if err := tx.Expiration.UnmarshalJSON([]byte(`"2016-04-06T08:29:27"`)); err != nil {
		suite.Fail(err.Error(), "Unmarshal expiration")
	}

	suite.RefTx = tx
}

func (suite *operationsAPITest) TearDown() {
	if err := suite.TestAPI.Close(); err != nil {
		suite.Fail(err.Error(), "Close")
	}
}

func (suite *operationsAPITest) Test_SerializeEmptyTransaction() {

	tx := types.NewTransaction()
	if err := tx.Expiration.UnmarshalJSON([]byte(`"2016-04-06T08:29:27"`)); err != nil {
		suite.Fail(err.Error(), "Unmarshal expiration")
	}

	suite.compareTransaction(tx)
}

func (suite *operationsAPITest) Test_SerializeTransaction() {
	hex, err := suite.TestAPI.SerializeTransaction(suite.RefTx)
	if err != nil {
		suite.Fail(err.Error(), "SerializeTransaction")
	}

	suite.NotNil(hex)
	suite.Equal("f68585abf4dce7c80457000000", hex)
}

// func (suite *operationsAPITest) Test_TransferOperation() {

// 	// 	wif1 := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
// 	// 	wif2 := "5KCBDTcyDqzsqehcb52tW5nU6pXife6V2rX9Yf7c3saYSzbDZ5W"
// 	// 	pub1 := crypto.GetPublicKey(wif1)
// 	// 	pub2 := crypto.GetPublicKey(wif2)

// 	// message, err := hex.DecodeString("abcdef0123456789")
// 	// if err != nil {
// 	// 	suite.Fail(err.Error(), "DecodeString")
// 	// }

// 	nonce := types.UInt64(5862723643998573708)

// 	op := types.TransferOperation{
// 		Extensions: types.Extensions{},
// 		Memo: types.Memo{
// 			From:    types.PublicKey(TestAccount3PubKeyActive),
// 			To:      types.PublicKey(TestAccount3PubKeyOwner),
// 			Nonce:   nonce,
// 			Message: "abcdef0123456789",
// 		},
// 		From: *types.NewGrapheneID("1.2.29"),
// 		To:   *types.NewGrapheneID("1.2.30"),
// 		Amount: types.AssetAmount{
// 			Amount: 10000,
// 			Asset:  *types.NewGrapheneID("1.3.5"),
// 		},
// 		Fee: types.AssetAmount{
// 			Amount: 256,
// 			Asset:  *types.NewGrapheneID("1.3.0"),
// 		},
// 	}

// 	suite.RefTx.Operations = types.Operations{
// 		types.Operation(&op),
// 	}

// 	suite.compareTransaction(suite.RefTx)
// }

// func (suite *operationsAPITest) Test_AccountCreateOperation() {
// 	op := types.AccountCreateOperation{
// 		Fee: types.AssetAmount{
// 			Amount: 123,
// 			Asset:  *types.NewGrapheneID("1.3.0"),
// 		},
// 		Registrar: *types.NewGrapheneID("1.2.345"),
// 		Referrer:  *types.NewGrapheneID("1.2.123"),
// 		Name:      "lala-account",
// 		Owner: types.Authority{
// 			WeightThreshold: 1,
// 			AccountAuths: types.AccountAuthsMap{
// 				*types.NewGrapheneID("1.2.4567"): 1,
// 				*types.NewGrapheneID("1.2.8904"): 2,
// 			},
// 			KeyAuths: types.KeyAuthsMap{
// 				types.PublicKey("BTS6zLNtyFVToBsBZDsgMhgjpwysYVbsQD6YhP3kRkQhANUB4w7Qp"): 1,
// 				//types.PublicKey("BTS6pbVDAjRFiw6fkiKYCrkz7PFeL7XNAfefrsREwg8MKpJ9VYV9x"): 2,
// 			},
// 			AddressAuths: types.AuthsMap{},
// 			Extensions:   types.Extensions{},
// 		},
// 		Active: types.Authority{
// 			WeightThreshold: 1,
// 			AccountAuths:    types.AccountAuthsMap{},
// 			KeyAuths: types.KeyAuthsMap{
// 				types.PublicKey("BTS6pbVDAjRFiw6fkiKYCrkz7PFeL7XNAfefrsREwg8MKpJ9VYV9x"): 1,
// 				//types.PublicKey("BTS6zLNtyFVToBsBZDsgMhgjpwysYVbsQD6YhP3kRkQhANUB4w7Qp"): 2,
// 				//types.PublicKey("BTS8CemMDjdUWSV5wKotEimhK6c4dY7p2PdzC2qM1HpAP8aLtZfE7"): 3,
// 			},
// 			AddressAuths: types.AuthsMap{},
// 			Extensions:   types.Extensions{},
// 		},
// 		Options: types.AccountOptions{
// 			MemoKey:       types.PublicKey("BTS5TPTziKkLexhVKsQKtSpo4bAv5RnB8oXcG4sMHEwCcTf3r7dqE"),
// 			VotingAccount: *types.NewGrapheneID("1.2.5"),
// 			NumWitness:    2,
// 			NumCommittee:  3,
// 			Votes: types.Votes{
// 				*types.NewVoteID("123:456"),
// 				*types.NewVoteID("789:123"),
// 			},
// 			Extensions: types.Extensions{},
// 		},
// 		Extensions: types.AccountCreateExtensions{
// 			BuybackOptions: types.BuybackOptions{
// 				AssetToBuy:       *types.NewGrapheneID("1.3.127"),
// 				AssetToBuyIssuer: *types.NewGrapheneID("1.2.31"),
// 				// Markets: types.GrapheneIDs{
// 				// 	*types.NewGrapheneID("1.3.20"),
// 				// 	*types.NewGrapheneID("1.3.21"),
// 				// 	*types.NewGrapheneID("1.3.22"),
// 				// },
// 			},
// 			//NullExtension: types.NullExtension{}{},
// 			//	OwnerSpecialAuthority:
// 		},
// 	}

// 	suite.RefTx.Operations = types.Operations{
// 		types.Operation(&op),
// 	}

// 	suite.compareTransaction(suite.RefTx)
// }

func (suite *operationsAPITest) Test_AssetReserveOperation() {
	op := AssetReserveOperation{
		Extensions: types.Extensions{},
		Fee: types.AssetAmount{
			Amount: 123,
			Asset:  *types.NewGrapheneID("1.3.0"),
		},

		Payer: *types.NewGrapheneID("1.2.25"),
		AmountToReserve: types.AssetAmount{
			Amount: 123456,
			Asset:  *types.NewGrapheneID("1.3.0"),
		},
	}

	suite.RefTx.Operations = types.Operations{
		types.Operation(&op),
	}

	suite.compareTransaction(suite.RefTx)
}

func (suite *operationsAPITest) Test_CallOrderUpdateOperation() {
	op := CallOrderUpdateOperation{
		Extensions:     types.Extensions{},
		FundingAccount: *types.NewGrapheneID("1.2.29"),
		DeltaCollateral: types.AssetAmount{
			Amount: 100000000,
			Asset:  *types.NewGrapheneID("1.3.0"),
		},
		DeltaDebt: types.AssetAmount{
			Amount: 10000,
			Asset:  *types.NewGrapheneID("1.3.22"),
		},
		Fee: types.AssetAmount{
			Amount: 100,
			Asset:  *types.NewGrapheneID("1.3.0"),
		},
	}

	suite.RefTx.Operations = types.Operations{
		types.Operation(&op),
	}

	suite.compareTransaction(suite.RefTx)
}

func (suite *operationsAPITest) Test_LimitOrderCreateOperation() {
	op := LimitOrderCreateOperation{
		Extensions: types.Extensions{},
		Seller:     *types.NewGrapheneID("1.2.29"),
		FillOrKill: false,
		Fee: types.AssetAmount{
			Amount: 100,
			Asset:  *types.NewGrapheneID("1.3.0"),
		},
		AmountToSell: types.AssetAmount{
			Amount: 100000,
			Asset:  *types.NewGrapheneID("1.3.0"),
		},
		MinToReceive: types.AssetAmount{
			Amount: 10000,
			Asset:  *types.NewGrapheneID("1.3.105"),
		},
	}

	if err := op.Expiration.UnmarshalJSON([]byte(`"2016-05-18T09:22:05"`)); err != nil {
		suite.Fail(err.Error(), "Unmarshal expiration")
	}

	suite.RefTx.Operations = types.Operations{
		types.Operation(&op),
	}

	suite.compareTransaction(suite.RefTx)
}

func (suite *operationsAPITest) Test_LimitOrderCancelOperation() {
	op := LimitOrderCancelOperation{
		Extensions:       types.Extensions{},
		Order:            *types.NewGrapheneID("1.7.123"),
		FeePayingAccount: *types.NewGrapheneID("1.2.456"),
		Fee: types.AssetAmount{
			Amount: 1000,
			Asset:  *types.NewGrapheneID("1.3.789"),
		},
	}

	suite.RefTx.Operations = types.Operations{
		types.Operation(&op),
	}

	suite.compareTransaction(suite.RefTx)
}

func (suite *operationsAPITest) compareTransaction(tx *types.Transaction) {
	var buf bytes.Buffer
	enc := util.NewTypeEncoder(&buf)
	if err := enc.Encode(tx); err != nil {
		suite.Fail(err.Error(), "Encode")
	}

	ref, err := suite.TestAPI.SerializeTransaction(tx)
	if err != nil {
		suite.Fail(err.Error(), "SerializeTransaction")
	}

	test := hex.EncodeToString(buf.Bytes())
	suite.Equal(ref, test)
}

func TestOperations(t *testing.T) {
	testSuite := new(operationsAPITest)
	suite.Run(t, testSuite)
	testSuite.TearDown()
}
