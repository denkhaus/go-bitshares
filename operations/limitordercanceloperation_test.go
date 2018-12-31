package operations

import (
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_LimitOrderCancelOperation() {

	op := LimitOrderCancelOperation{
		Extensions: types.Extensions{},
	}

	suite.OpSamplesTest(&op)
}
