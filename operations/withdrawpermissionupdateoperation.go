package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeWithdrawPermissionUpdate] = func() types.Operation {
		op := &WithdrawPermissionUpdateOperation{}
		return op
	}
}

type WithdrawPermissionUpdateOperation struct {
	types.OperationFee
	WithdrawFromAccount    types.AccountID            `json:"withdraw_from_account"`
	AuthorizedAccount      types.AccountID            `json:"authorized_account"`
	PermissionToUpdate     types.WithdrawPermissionID `json:"permission_to_update"`
	WithdrawalLimit        types.AssetAmount          `json:"withdrawal_limit"`
	WithdrawalPeriodSec    types.UInt32               `json:"withdrawal_period_sec"`
	PeriodStartTime        types.Time                 `json:"period_start_time"`
	PeriodsUntilExpiration types.UInt32               `json:"periods_until_expiration"`
}

func (p WithdrawPermissionUpdateOperation) Type() types.OperationType {
	return types.OperationTypeWithdrawPermissionUpdate
}

func (p WithdrawPermissionUpdateOperation) Marshal(enc *util.TypeEncoder) error {
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
	if err := enc.Encode(p.PermissionToUpdate); err != nil {
		return errors.Annotate(err, "encode PermissionToUpdate")
	}
	if err := enc.Encode(p.WithdrawalLimit); err != nil {
		return errors.Annotate(err, "encode WithdrawalLimit")
	}
	if err := enc.Encode(p.WithdrawalPeriodSec); err != nil {
		return errors.Annotate(err, "encode WithdrawalPeriodSec")
	}
	if err := enc.Encode(p.PeriodStartTime); err != nil {
		return errors.Annotate(err, "encode PeriodStartTime")
	}
	if err := enc.Encode(p.PeriodsUntilExpiration); err != nil {
		return errors.Annotate(err, "encode PeriodsUntilExpiration")
	}

	return nil
}
