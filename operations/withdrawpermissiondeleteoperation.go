package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeWithdrawPermissionDelete] = func() types.Operation {
		op := &WithdrawPermissionDeleteOperation{}
		return op
	}
}

type WithdrawPermissionDeleteOperation struct {
	types.OperationFee
	AuthorizedAccount    types.AccountID            `json:"authorized_account"`
	WithdrawFromAccount  types.AccountID            `json:"withdraw_from_account"`
	WithdrawalPermission types.WithdrawPermissionID `json:"withdrawal_permission"`
}

func (p WithdrawPermissionDeleteOperation) Type() types.OperationType {
	return types.OperationTypeWithdrawPermissionDelete
}

func (p WithdrawPermissionDeleteOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}
	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}
	if err := enc.Encode(p.WithdrawFromAccount); err != nil {
		return errors.Annotate(err, "encode WithdrawFromAccount")
	}
	if err := enc.Encode(p.AuthorizedAccount); err != nil {
		return errors.Annotate(err, "encode AuthorizedAccount")
	}
	if err := enc.Encode(p.WithdrawalPermission); err != nil {
		return errors.Annotate(err, "encode WithdrawalPermission")
	}

	return nil
}
