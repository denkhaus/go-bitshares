package operations

import (
	"github.com/denkhaus/bitshares/types"
)

func (suite *operationsAPITest) Test_LimitOrderCreateOperation() {
	op := LimitOrderCreateOperation{
		Extensions: types.Extensions{},
	}

	suite.OpSamplesTest(&op)
}
