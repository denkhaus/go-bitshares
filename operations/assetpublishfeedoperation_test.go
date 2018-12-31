package operations

import (
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_AssetPublishFeedOperation() {
	op := AssetPublishFeedOperation{
		Extensions: types.Extensions{},
	}

	suite.OpSamplesTest(&op)
}
