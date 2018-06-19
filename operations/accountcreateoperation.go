package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &AccountCreateOperation{}
	types.OperationMap[op.Type()] = op
}

type AccountCreateOperation struct {
	Registrar       types.GrapheneID              `json:"registrar"`
	Referrer        types.GrapheneID              `json:"referrer"`
	ReferrerPercent types.UInt16                  `json:"referrer_percent"`
	Owner           types.Authority               `json:"owner"`
	Active          types.Authority               `json:"active"`
	Fee             types.AssetAmount             `json:"fee"`
	Name            string                        `json:"name"`
	Extensions      types.AccountCreateExtensions `json:"extensions"`
	Options         types.AccountOptions          `json:"options"`
}

func (p AccountCreateOperation) GetFee() types.AssetAmount {
	return p.Fee
}

func (p *AccountCreateOperation) SetFee(fee types.AssetAmount) {
	p.Fee = fee
}

func (p AccountCreateOperation) Type() types.OperationType {
	return types.OperationTypeAccountCreate
}

func (p AccountCreateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
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

	if err := enc.Encode(p.Name); err != nil {
		return errors.Annotate(err, "encode name")
	}

	if err := enc.Encode(p.Owner); err != nil {
		return errors.Annotate(err, "encode owner")
	}

	if err := enc.Encode(p.Active); err != nil {
		return errors.Annotate(err, "encode active")
	}

	if err := enc.Encode(p.Options); err != nil {
		return errors.Annotate(err, "encode options")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

//NewAccountCreateOperation creates a new AccountCreateOperation
func NewAccountCreateOperation() *AccountCreateOperation {
	tx := AccountCreateOperation{
		Extensions: types.AccountCreateExtensions{},
	}
	return &tx
}
