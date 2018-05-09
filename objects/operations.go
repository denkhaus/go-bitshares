package objects

import (
	"encoding/json"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

var (
	opMap = make(map[OperationType]Operation)
)

type Operation interface {
	util.TypeMarshaller
	ApplyFee(fee AssetAmount)
	Type() OperationType
}

type OperationResult interface {
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

func (p *OperationEnvelope) UnmarshalJSON(data []byte) error {
	raw := make([]json.RawMessage, 2)
	if err := json.Unmarshal(data, &raw); err != nil {
		return errors.Annotate(err, "unmarshal raw object")
	}

	if len(raw) != 2 {
		return errors.Errorf("Invalid operation data: %v", string(data))
	}

	if err := json.Unmarshal(raw[0], &p.Type); err != nil {
		return errors.Annotate(err, "unmarshal OperationType")
	}

	if op, ok := opMap[p.Type]; ok {
		p.Operation = op
		if err := ffjson.Unmarshal(raw[1], p.Operation); err != nil {
			return errors.Annotate(err, "unmarshal Operation")
		}
	} else {
		util.DumpUnmarshaled("Operation", raw[1])
		return errors.Errorf("Operation type %d not yet supported", p.Type)
	}

	return nil
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

func (p *Operations) UnmarshalJSON(data []byte) error {
	var envs []OperationEnvelope
	if err := json.Unmarshal(data, &envs); err != nil {
		return err
	}

	ops := make(Operations, len(envs))
	for idx, env := range envs {
		ops[idx] = env.Operation.(Operation)
	}

	*p = ops
	return nil
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
