package data

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/juju/errors"
)

type OperationSampleMap map[types.OperationType]string

var (
	OpSampleMap              = make(OperationSampleMap)
	ErrNoSampleDataAvailable = errors.New("no sample data available")
)

//GetSampleByType returns a Operation data sample by OperationID
func GetSampleByType(typ types.OperationType) (string, error) {
	if s, ok := OpSampleMap[typ]; ok {
		return s, nil
	}

	return "", ErrNoSampleDataAvailable
}
