package operations

import (
	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/types"
)

// op := AssetReserveOperation{
// 	Extensions: types.Extensions{},
// 	Fee: types.AssetAmount{
// 		Amount: 123,
// 		Asset:  *types.NewGrapheneID("1.3.0"),
// 	},

// 	Payer: *types.NewGrapheneID("1.2.25"),
// 	AmountToReserve: types.AssetAmount{
// 		Amount: 123456,
// 		Asset:  *types.NewGrapheneID("1.3.0"),
// 	},
// }

func (suite *operationsAPITest) Test_AssetReserveOperation() {
	op := AssetReserveOperation{
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
