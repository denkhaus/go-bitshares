package objects

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &AccountCreateOperation{}
	opMap[op.Type()] = op
}

type AccountCreateOperation struct {
	Registrar       GrapheneID  `json:"registrar"`
	Referrer        GrapheneID  `json:"referrer"`
	ReferrerPercent UInt32      `json:"referrer_percent"`
	Owner           Authority   `json:"owner"`
	Active          Authority   `json:"active"`
	Fee             AssetAmount `json:"fee"`
	Name            string      `json:"name"`
	//	Extensions      Extensions     `json:"extensions"`
	Options AccountOptions `json:"options"`
}

//implements Operation interface
func (p *AccountCreateOperation) ApplyFee(fee AssetAmount) {
	p.Fee = fee
}

//implements Operation interface
func (p AccountCreateOperation) Type() OperationType {
	return OperationTypeAccountCreate
}

//TODO: validate encode order!
//implements Operation interface
func (p AccountCreateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Registrar); err != nil {
		return errors.Annotate(err, "encode registrar")
	}

	if err := enc.Encode(p.Referrer); err != nil {
		return errors.Annotate(err, "encode referrer")
	}

	if err := enc.Encode(p.ReferrerPercent); err != nil {
		return errors.Annotate(err, "encode referrer percent")
	}

	if err := enc.Encode(p.Owner); err != nil {
		return errors.Annotate(err, "encode owner")
	}

	if err := enc.Encode(p.Active); err != nil {
		return errors.Annotate(err, "encode active")
	}

	if err := enc.Encode(p.Name); err != nil {
		return errors.Annotate(err, "encode name")
	}

	// if err := enc.Encode(p.Extensions); err != nil {
	// 	return errors.Annotate(err, "encode extensions")
	// }

	if err := enc.Encode(p.Options); err != nil {
		return errors.Annotate(err, "encode options")
	}
	return nil
}

//NewAccountCreateOperation creates a new AccountCreateOperation
func NewAccountCreateOperation() *AccountCreateOperation {
	tx := AccountCreateOperation{
	//	Extensions: Extensions{},
	}
	return &tx
}
