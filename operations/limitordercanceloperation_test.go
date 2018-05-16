package operations

import (
	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/types"
)

// op := LimitOrderCancelOperation{
// 	Extensions:       types.Extensions{},
// 	Order:            *types.NewGrapheneID("1.7.123"),
// 	FeePayingAccount: *types.NewGrapheneID("1.2.456"),
// 	Fee: types.AssetAmount{
// 		Amount: 1000,
// 		Asset:  *types.NewGrapheneID("1.3.789"),
// 	},
// }

func (suite *operationsAPITest) Test_LimitOrderCancelOperation() {

	op := LimitOrderCancelOperation{
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
