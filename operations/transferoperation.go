package operations

//go:generate ffjson   $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &TransferOperation{}
	types.OperationMap[op.Type()] = op
}

type TransferOperation struct {
	From       types.GrapheneID  `json:"from"`
	To         types.GrapheneID  `json:"to"`
	Amount     types.AssetAmount `json:"amount"`
	Fee        types.AssetAmount `json:"fee"`
	Memo       *types.Memo       `json:"memo,omitempty"`
	Extensions types.Extensions  `json:"extensions"`
}

func (p *TransferOperation) ApplyFee(fee types.AssetAmount) {
	p.Fee = fee
}

func (p TransferOperation) Type() types.OperationType {
	return types.OperationTypeTransfer
}

func (p TransferOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.From); err != nil {
		return errors.Annotate(err, "encode from")
	}

	if err := enc.Encode(p.To); err != nil {
		return errors.Annotate(err, "encode to")
	}

	if err := enc.Encode(p.Amount); err != nil {
		return errors.Annotate(err, "encode amount")
	}

	if err := enc.Encode(p.Memo != nil); err != nil {
		return errors.Annotate(err, "encode have memo")
	}

	if err := enc.Encode(p.Memo); err != nil {
		return errors.Annotate(err, "encode memo")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

//NewTransferOperation creates a new TransferOperation
func NewTransferOperation() *TransferOperation {
	tx := TransferOperation{
		Extensions: types.Extensions{},
	}
	return &tx
}
