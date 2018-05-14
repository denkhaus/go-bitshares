package main

import (
	"fmt"
	"reflect"

	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/types"
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

type OpDataStore map[types.OperationType]*OperationBlob

//TODO: save last scanned block and reapply
func (p *OpDataStore) Init(m data.OperationSampleMap, ch chan GenData) error {
	if len(m) == 0 {
		fmt.Printf("init datastore: no sample data loaded\n")
		return errors.New("no sample data loaded")
	}

	for typ, data := range m {
		opData, err := objx.FromJSON(data)
		if err != nil {
			return errors.Annotate(err, "FromJSON")
		}

		ok, err := p.Evaluate(typ, NewOperationBlob(opData), 0)
		if err != nil {
			return errors.Annotate(err, "Evaluate")
		}

		if ok {
			ch <- GenData{
				Type: typ,
				Data: opData,
			}
		}
	}

	return nil
}

func (p *OpDataStore) Evaluate(typ types.OperationType, blob *OperationBlob, block uint64) (bool, error) {
	if err := reflectwalk.Walk(blob.data, blob); err != nil {
		return false, errors.Annotate(err, "Walk")
	}

	if bl, ok := (*p)[typ]; ok {
		if bl.fields < blob.fields {
			(*p)[typ] = blob
			fmt.Printf("Block %d: Blob upgraded for type %v: %s\n", block, typ, blob)
			return true, nil
		}
	} else {
		(*p)[typ] = blob
		fmt.Printf("Block %d: Blob added for type %v: %s\n", block, typ, blob)
		return true, nil
	}

	return false, nil
}

func NewOpDataStore() *OpDataStore {
	s := OpDataStore{}
	return &s
}
