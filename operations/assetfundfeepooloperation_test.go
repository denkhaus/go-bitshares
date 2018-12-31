package operations

import (
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_AssetFundFeePoolOperation() {
	op := AssetFundFeePoolOperation{
		Extensions: types.Extensions{},
	}

	suite.OpSamplesTest(&op)
}
