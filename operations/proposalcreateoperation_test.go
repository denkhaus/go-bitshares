package operations

import (
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_ProposalCreateOperation() {
	op := ProposalCreateOperation{
		Extensions: types.Extensions{},
	}

	suite.OpSamplesTest(&op)
}
