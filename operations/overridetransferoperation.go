package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &OverrideTransferOperation{}
	types.OperationMap[op.Type()] = op
}

type OverrideTransferOperation struct {
	Amount     types.AssetAmount `json:"amount"`
	Extensions types.Extensions  `json:"extensions"`
	Fee        types.AssetAmount `json:"fee"`
	From       types.GrapheneID  `json:"from"`
	Issuer     types.GrapheneID  `json:"issuer"`
	Memo       *types.Memo       `json:"memo,omitempty"`
	To         types.GrapheneID  `json:"to"`
}

func (p OverrideTransferOperation) GetFee() types.AssetAmount {
	return p.Fee
}

func (p *OverrideTransferOperation) SetFee(fee types.AssetAmount) {
	p.Fee = fee
}

func (p OverrideTransferOperation) Type() types.OperationType {
	return types.OperationTypeOverrideTransfer
}

func (p OverrideTransferOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Issuer); err != nil {
		return errors.Annotate(err, "encode Issuer")
	}

	if err := enc.Encode(p.From); err != nil {
		return errors.Annotate(err, "encode From")
	}

	if err := enc.Encode(p.To); err != nil {
		return errors.Annotate(err, "encode To")
	}

	if err := enc.Encode(p.Amount); err != nil {
		return errors.Annotate(err, "encode Amount")
	}

	if err := enc.Encode(p.Memo != nil); err != nil {
		return errors.Annotate(err, "encode have Memo")
	}

	if err := enc.Encode(p.Memo); err != nil {
		return errors.Annotate(err, "encode Memo")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

//NewOverrideTransferOperation creates a new OverrideTransferOperation
func NewOverrideTransferOperation() *OverrideTransferOperation {
	tx := OverrideTransferOperation{
		Extensions: types.Extensions{},
	}
	return &tx
}
