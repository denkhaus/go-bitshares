package operations

import (
	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_AssetPublishFeedOperation() {
	op := AssetPublishFeedOperation{
		Extensions: types.Extensions{},
	}

	samples, err := data.GetSamplesByType(op.Type())
	if err != nil {
		suite.FailNow(err.Error(), "GetSamplesByType")
	}

	for idx, sample := range samples {
		if err := op.UnmarshalJSON([]byte(sample)); err != nil {
			suite.FailNow(err.Error(), "UnmarshalJSON")
		}

		suite.RefTx.Operations = types.Operations{
			types.Operation(&op),
		}

		suite.compareTransaction(idx, suite.RefTx, false)
	}
}
