package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeWithdrawPermissionCreate] = func() types.Operation {
		op := &WithdrawPermissionCreateOperation{}
		return op
	}
}

type WithdrawPermissionCreateOperation struct {
	types.OperationFee
	AuthorizedAccount      types.GrapheneID  `json:"authorized_account"`
	PeriodStartTime        types.Time        `json:"period_start_time"`
	PeriodsUntilExpiration types.UInt32      `json:"periods_until_expiration"`
	WithdrawFromAccount    types.GrapheneID  `json:"withdraw_from_account"`
	WithdrawalLimit        types.AssetAmount `json:"withdrawal_limit"`
	WithdrawalPeriodSec    types.UInt32      `json:"withdrawal_period_sec"`
}

func (p WithdrawPermissionCreateOperation) Type() types.OperationType {
	return types.OperationTypeWithdrawPermissionCreate
}

func (p WithdrawPermissionCreateOperation) Marshal(enc *util.TypeEncoder) error {
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
	if err := enc.Encode(p.WithdrawalLimit); err != nil {
		return errors.Annotate(err, "encode WithdrawalLimit")
	}
	if err := enc.Encode(p.WithdrawalPeriodSec); err != nil {
		return errors.Annotate(err, "encode WithdrawalPeriodSec")
	}
	if err := enc.Encode(p.PeriodsUntilExpiration); err != nil {
		return errors.Annotate(err, "encode PeriodsUntilExpiration")
	}
	if err := enc.Encode(p.PeriodStartTime); err != nil {
		return errors.Annotate(err, "encode PeriodStartTime")
	}

	return nil
}
