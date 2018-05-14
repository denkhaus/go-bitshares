package data

import (
	"github.com/denkhaus/bitshares/types"
)

type OperationSampleMap map[types.OperationType]string

var (
	OpSampleMap = make(OperationSampleMap)
)

//GetSampleByType returns a Operation data sample by OperationID
func GetSampleByType(typ types.OperationType) string {
	if s, ok := OpSampleMap[typ]; ok {
		return s
	}

	return ""
}
