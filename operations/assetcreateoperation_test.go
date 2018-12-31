package operations

import (
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_AssetCreateOperation() {
	op := AssetCreateOperation{
		Extensions: types.Extensions{},
	}

	suite.OpSamplesTest(&op)
}
