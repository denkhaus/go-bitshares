package operations

import (
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_AssetUpdateFeedProducersOperation() {
	op := AssetUpdateFeedProducersOperation{
		Extensions: types.Extensions{},
	}

	suite.OpSamplesTest(&op)
}
