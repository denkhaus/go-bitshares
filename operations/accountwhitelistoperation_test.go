package operations

import (
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_AccountWhitelistOperation() {
	op := AccountWhitelistOperation{
		Extensions: types.Extensions{},
	}

	suite.OpSamplesTest(&op)
}
