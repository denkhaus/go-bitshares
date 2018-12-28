package data

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/juju/errors"
)

type OperationSampleMap map[types.OperationType]map[int]string

var (
	OpSampleMap              OperationSampleMap
	ErrNoSampleDataAvailable = errors.New("no sample data available")
)

func init() {
	OpSampleMap = make(OperationSampleMap)
}

//GetSampleByType returns a Operation data sample by OperationID
func GetSamplesByType(typ types.OperationType) (map[int]string, error) {
	ret := make(map[int]string)
	if s, ok := OpSampleMap[typ]; ok {
		return s, nil
	}

	return ret, ErrNoSampleDataAvailable
}
