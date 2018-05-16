package operations

import (
	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/types"
)

// op := LimitOrderCreateOperation{
// 	Extensions: types.Extensions{},
// 	Seller:     *types.NewGrapheneID("1.2.29"),
// 	FillOrKill: false,
// 	Fee: types.AssetAmount{
// 		Amount: 100,
// 		Asset:  *types.NewGrapheneID("1.3.0"),
// 	},
// 	AmountToSell: types.AssetAmount{
// 		Amount: 100000,
// 		Asset:  *types.NewGrapheneID("1.3.0"),
// 	},
// 	MinToReceive: types.AssetAmount{
// 		Amount: 10000,
// 		Asset:  *types.NewGrapheneID("1.3.105"),
// 	},
// }
func (suite *operationsAPITest) Test_LimitOrderCreateOperation() {
	op := LimitOrderCreateOperation{
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
