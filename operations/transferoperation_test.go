package operations

import (
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_TransferOperation() {
	op := TransferOperation{
		Extensions: types.Extensions{},
	}

	suite.OpSamplesTest(&op)
}
