package operations

import (
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_OverrideTransferOperation() {
	op := OverrideTransferOperation{
		Extensions: types.Extensions{},
	}

	suite.OpSamplesTest(&op)
}
