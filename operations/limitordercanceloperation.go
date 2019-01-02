package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeLimitOrderCancel] = func() types.Operation {
		op := &LimitOrderCancelOperation{}
		return op
	}
}

type LimitOrderCancelOperation struct {
	types.OperationFee
	FeePayingAccount types.AccountID    `json:"fee_paying_account"`
	Order            types.LimitOrderID `json:"order"`
	Extensions       types.Extensions   `json:"extensions"`
}

func (p LimitOrderCancelOperation) Type() types.OperationType {
	return types.OperationTypeLimitOrderCancel
}

func (p LimitOrderCancelOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}
	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}
	if err := enc.Encode(p.FeePayingAccount); err != nil {
		return errors.Annotate(err, "encode FeePayingAccount")
	}
	if err := enc.Encode(p.Order); err != nil {
		return errors.Annotate(err, "encode Order")
	}
	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extensions")
	}

	return nil
}
