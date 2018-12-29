package main

import (
	"fmt"

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
	data map[string]interface{}
	rank int
}

func (p *OperationBlob) Enter(loc reflectwalk.Location) error {
	p.rank++
	return nil
}
func (p *OperationBlob) Exit(loc reflectwalk.Location) error {
	return nil
}

func (p *OperationBlob) String() string {
	return fmt.Sprintf("rank: %d", p.rank)
}

func (p *OperationBlob) Analyze() error {
	if err := reflectwalk.Walk(p.data, p); err != nil {
		return errors.Annotate(err, "Walk")
	}

	return nil
}

func NewOperationBlob(data map[string]interface{}) *OperationBlob {
	s := OperationBlob{
		data: data,
	}
	return &s
}

type OpDataStore map[types.OperationType]map[int]*OperationBlob

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

func (p *OpDataStore) MustInsert(typ types.OperationType, rank int) int {
	blobs, ok := (*p)[typ]
	if ok {
		for _, bl := range blobs {
			if bl.rank == rank {
				return -1
			}
		}

		return len(blobs)
	}

	return 0
}

func (p *OpDataStore) Insert(typ types.OperationType, blob *OperationBlob, block uint64) (int, error) {
	if err := blob.Analyze(); err != nil {
		return -1, errors.Annotate(err, "Analyze")
	}

	idx := p.MustInsert(typ, blob.rank)
	if idx < 0 {
		return idx, nil
	}

	if idx == 0 {
		(*p)[typ] = map[int]*OperationBlob{0: blob}
		logging.Infof("Block %d: Blob added for type %v: %s", block, typ, blob)
		return idx, nil
	}

	(*p)[typ][idx] = blob
	logging.Infof("Block %d: Blob upgraded for type %v: %s", block, typ, blob)
	return idx, nil
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
