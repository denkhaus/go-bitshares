package operations

import (
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_AccountUpgradeOperation() {
	op := AccountUpgradeOperation{
		Extensions: types.Extensions{},
	}

	suite.OpSamplesTest(&op)
}
