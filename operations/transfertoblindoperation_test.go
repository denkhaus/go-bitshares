package operations

// func (suite *operationsAPITest) Test_TransferToBlindOperation() {
// 	op := TransferToBlindOperation{}

// 	sample, err := data.GetSampleByType(op.Type())
// 	if err != nil {
// 		suite.FailNow(err.Error(), "GetSampleByType")
// 	}

// 	if err := op.UnmarshalJSON([]byte(sample)); err != nil {
// 		suite.FailNow(err.Error(), "UnmarshalJSON")
// 	}

// 	suite.RefTx.Operations = types.Operations{
// 		types.Operation(&op),
// 	}

// 	suite.compareTransaction(suite.RefTx, false)
// }
