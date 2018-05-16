package operations

import (
	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_TransferOperation() {
	op := TransferOperation{
		Extensions: types.Extensions{},
	}

	sample, err := data.GetSampleByType(op.Type())
	if err != nil {
		suite.FailNow(err.Error(), "GetSampleByType")
	}

	if err := op.UnmarshalJSON([]byte(sample)); err != nil {
		suite.FailNow(err.Error(), "UnmarshalJSON")
	}

	suite.RefTx.Operations = types.Operations{
		types.Operation(&op),
	}

	suite.compareTransaction(suite.RefTx)
}

// func (suite *operationsAPITest) Test_TransferOperation() {

// 	// 	wif1 := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
// 	// 	wif2 := "5KCBDTcyDqzsqehcb52tW5nU6pXife6V2rX9Yf7c3saYSzbDZ5W"
// 	// 	pub1 := crypto.GetPublicKey(wif1)
// 	// 	pub2 := crypto.GetPublicKey(wif2)

// 	// message, err := hex.DecodeString("abcdef0123456789")
// 	// if err != nil {
// 	// 	suite.FailNow(err.Error(), "DecodeString")
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
