package data

import "github.com/denkhaus/bitshares/types"

type OperationSampleMap map[types.OperationType]string

var (
	OpSampleMap = make(OperationSampleMap)
)
