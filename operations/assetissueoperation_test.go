package operations

import (
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_AssetIssueOperation() {
	op := AssetIssueOperation{
		Extensions: types.Extensions{},
	}

	suite.OpSamplesTest(&op)
}
