package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/denkhaus/logging"
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

	// asset                         fee;
	// account_id_type               withdraw_from_account;
	// account_id_type               authorized_account;
	// withdraw_permission_id_type   permission_to_update;
	// asset                         withdrawal_limit;
	// uint32_t                      withdrawal_period_sec = 0;
	// time_point_sec                period_start_time;
	// uint32_t                      periods_until_expiration = 0;
}

func (p WithdrawPermissionUpdateOperation) Type() types.OperationType {
	return types.OperationTypeWithdrawPermissionCreate
}

func (p WithdrawPermissionUpdateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}
	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}

	// (fee)(withdraw_from_account)(authorized_account)
	//   (permission_to_update)(withdrawal_limit)(withdrawal_period_sec)(period_start_time)(periods_until_expiration)

	logging.Warnf("%s is not implemented", p.Type().OperationName())
	return nil
}
