package operations

import (
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_AssetUpdateBitassetOperation() {
	op := AssetUpdateBitassetOperation{
		Extensions: types.Extensions{},
	}

	suite.OpSamplesTest(&op)
}
