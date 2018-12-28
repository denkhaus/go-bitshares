package main

import (
	"fmt"
	"reflect"

	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/logging"
	"github.com/juju/errors"
	"github.com/mitchellh/reflectwalk"
	"github.com/stretchr/objx"
)

//OperationBlob simply counts fields and child map elements.
//The more fields, the more telling is the structure.
type OperationBlob struct {
	data   map[string]interface{}
	fields int
}

func (p *OperationBlob) Map(m reflect.Value) error {
	p.fields++
	return nil
}
func (p *OperationBlob) MapElem(m, k, v reflect.Value) error {
	p.fields++
	return nil
}

func (p *OperationBlob) String() string {
	return fmt.Sprintf("rank: %d", p.fields)
}

func NewOperationBlob(data map[string]interface{}) *OperationBlob {
	s := OperationBlob{
		data: data,
	}
	return &s
}

type OpDataStore map[types.OperationType]map[int]*OperationBlob

//TODO: save last scanned block and reapply
func (p *OpDataStore) Init(m data.OperationSampleMap, ch chan GenData) error {
	if len(m) == 0 {
		logging.Warn("init datastore: no sample data loaded")
		return nil
	}

	logging.Info("init datastore")
	for typ, d := range m {
		for idx, data := range d {
			opData, err := objx.FromJSON(data)
			if err != nil {
				return errors.Annotate(err, "FromJSON")
			}

			_, err = p.Insert(typ, NewOperationBlob(opData), 0)
			if err != nil {
				return errors.Annotate(err, "Insert")
			}

			if idx >= 0 {
				ch <- GenData{
					Type:      typ,
					SampleIdx: idx,
					Data:      opData,
				}
			}
		}
	}

	return nil
}

func (p *OpDataStore) Insert(typ types.OperationType, blob *OperationBlob, block uint64) (int, error) {
	if err := reflectwalk.Walk(blob.data, blob); err != nil {
		return -1, errors.Annotate(err, "Walk")
	}

	if blobs, ok := (*p)[typ]; ok {
		var maxFields = 0
		var maxIdx = -1
		for idx, bl := range blobs {
			maxFields = max(bl.fields, maxFields)
			maxIdx = max(idx, maxIdx)
		}

		if maxFields < blob.fields {
			newIdx := maxIdx + 1
			(*p)[typ][newIdx] = blob
			logging.Infof("Block %d: Blob upgraded for type %v: %s", block, typ, blob)
			return newIdx, nil
		}

	} else {
		(*p)[typ] = map[int]*OperationBlob{
			0: blob,
		}

		logging.Infof("Block %d: Blob added for type %v: %s", block, typ, blob)
		return 0, nil
	}

	return -1, nil
}

func NewOpDataStore() *OpDataStore {
	s := OpDataStore{}
	return &s
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
