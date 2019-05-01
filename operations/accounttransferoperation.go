package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeAccountTransfer] = func() types.Operation {
		op := &AccountTransferOperation{}
		return op
	}
}

type AccountTransferOperation struct {
	types.OperationFee
	AccountID  types.AccountID  `json:"account_id"`
	NewOwner   types.AccountID  `json:"new_owner"`
	Extensions types.Extensions `json:"extensions"`
}

func (p AccountTransferOperation) Type() types.OperationType {
	return types.OperationTypeAccountTransfer
}

func (p AccountTransferOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}
	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}
	if err := enc.Encode(p.AccountID); err != nil {
		return errors.Annotate(err, "encode AccountID")
	}
	if err := enc.Encode(p.NewOwner); err != nil {
		return errors.Annotate(err, "encode NewOwner")
	}
	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extensions")
	}
	return nil
}
