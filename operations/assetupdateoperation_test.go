package operations

import (
	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_AssetUpdateOperation() {
	op := AssetUpdateOperation{
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

	suite.compareTransaction(suite.RefTx, false)
}
