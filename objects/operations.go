package objects

import (
	"encoding/json"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type Operation interface {
	util.TypeMarshaller
	ApplyFee(fee AssetAmount)
	Type() OperationType
}

type OperationEnvelope struct {
	Type      OperationType
	Operation Operation
}

func (p OperationEnvelope) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		p.Type,
		p.Operation,
	})
}

type Operations []Operation

//implements TypeMarshaller interface
func (p Operations) Marshal(enc *util.TypeEncoder) error {

	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode Operations length")
	}

	for _, op := range p {
		if err := enc.Encode(op); err != nil {
			return errors.Annotate(err, "encode Operation")
		}
	}

	return nil
}

func (p Operations) MarshalJSON() ([]byte, error) {
	env := make([]OperationEnvelope, len(p))

	for idx, op := range p {
		env[idx] = OperationEnvelope{
			Type:      op.Type(),
			Operation: op,
		}
	}

	return json.Marshal(env)
}

func (p Operations) ApplyFees(fees []AssetAmount) error {
	if len(p) != len(fees) {
		return errors.New("count of fees must match count of operations")
	}

	for idx, op := range p {
		op.ApplyFee(fees[idx])
	}

	return nil
}

func (p Operations) Types() [][]OperationType {
	ret := make([][]OperationType, len(p))
	for idx, op := range p {
		ret[idx] = []OperationType{op.Type()}
	}

	return ret
}
