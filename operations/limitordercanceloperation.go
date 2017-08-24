package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/objects"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type LimitOrderCancelOperation struct {
	FeePayingAccount objects.GrapheneID  `json:"fee_paying_account"`
	Order            objects.GrapheneID  `json:"order"`
	Fee              objects.AssetAmount `json:"fee"`
	Extensions       objects.Extensions  `json:"extensions"`
}

//implements Operation interface
func (p *LimitOrderCancelOperation) ApplyFee(fee objects.AssetAmount) {
	p.Fee = fee
}

//implements Operation interface
func (p LimitOrderCancelOperation) Type() objects.OperationType {
	return objects.OperationTypeLimitOrderCancel
}

//implements Operation interface
func (p LimitOrderCancelOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
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

	/* if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	} */
	return nil
}

func NewLimitOrderCancelOperation(order objects.GrapheneObject) *LimitOrderCancelOperation {
	op := LimitOrderCancelOperation{
		Extensions: objects.Extensions{},
	}

	op.Order.FromObjectID(order.Id())
	return &op
}
