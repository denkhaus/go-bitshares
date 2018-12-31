package operations

import (
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_AssetUpdateOperation() {
	op := AssetUpdateOperation{
		Extensions: types.Extensions{},
	}

	suite.OpSamplesTest(&op)
}
