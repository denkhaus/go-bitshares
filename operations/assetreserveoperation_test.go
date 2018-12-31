package operations

import (
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_AssetReserveOperation() {
	op := AssetReserveOperation{
		Extensions: types.Extensions{},
	}

	suite.OpSamplesTest(&op)
}
