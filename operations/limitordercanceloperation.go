package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &LimitOrderCancelOperation{}
	types.OperationMap[op.Type()] = op
}

type LimitOrderCancelOperation struct {
	FeePayingAccount types.GrapheneID  `json:"fee_paying_account"`
	Order            types.GrapheneID  `json:"order"`
	Fee              types.AssetAmount `json:"fee"`
	Extensions       types.Extensions  `json:"extensions"`
}

func (p *LimitOrderCancelOperation) ApplyFee(fee types.AssetAmount) {
	p.Fee = fee
}

func (p LimitOrderCancelOperation) Type() types.OperationType {
	return types.OperationTypeLimitOrderCancel
}

func (p LimitOrderCancelOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.FeePayingAccount); err != nil {
		return errors.Annotate(err, "encode from")
	}

	if err := enc.Encode(p.Order); err != nil {
		return errors.Annotate(err, "encode to")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

func NewLimitOrderCancelOperation(order types.GrapheneID) *LimitOrderCancelOperation {
	op := LimitOrderCancelOperation{
		Extensions: types.Extensions{},
		Order:      order,
	}

	return &op
}
