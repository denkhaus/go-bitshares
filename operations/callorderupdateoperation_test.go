package operations

import (
	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/types"
)

// op := CallOrderUpdateOperation{
// 	Extensions:     types.Extensions{},
// 	FundingAccount: *types.NewGrapheneID("1.2.29"),
// 	DeltaCollateral: types.AssetAmount{
// 		Amount: 100000000,
// 		Asset:  *types.NewGrapheneID("1.3.0"),
// 	},
// 	DeltaDebt: types.AssetAmount{
// 		Amount: 10000,
// 		Asset:  *types.NewGrapheneID("1.3.22"),
// 	},
// 	Fee: types.AssetAmount{
// 		Amount: 100,
// 		Asset:  *types.NewGrapheneID("1.3.0"),
// 	},
// }

func (suite *operationsAPITest) Test_CallOrderUpdateOperation() {
	op := CallOrderUpdateOperation{
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
