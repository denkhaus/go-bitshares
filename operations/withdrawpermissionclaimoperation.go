package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/denkhaus/logging"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeWithdrawPermissionClaim] = func() types.Operation {
		op := &WithdrawPermissionClaimOperation{}
		return op
	}
}

type WithdrawPermissionClaimOperation struct {
	types.OperationFee

	// asset                       fee;
	// withdraw_permission_id_type withdraw_permission;
	// account_id_type             withdraw_from_account;
	// account_id_type             withdraw_to_account;
	// asset                       amount_to_withdraw;
	// optional<memo_data> memo;

}

func (p WithdrawPermissionClaimOperation) Type() types.OperationType {
	return types.OperationTypeWithdrawPermissionClaim
}

func (p WithdrawPermissionClaimOperation) MarshalFeeScheduleParams(params types.M, enc *util.TypeEncoder) error {
	if fee, ok := params["fee"]; ok {
		if err := enc.Encode(types.UInt64(fee.(float64))); err != nil {
			return errors.Annotate(err, "encode Fee")
		}
	}

	if ppk, ok := params["price_per_kbyte"]; ok {
		if err := enc.Encode(types.UInt32(ppk.(float64))); err != nil {
			return errors.Annotate(err, "encode PricePerKByte")
		}
	}

	return nil
}

func (p WithdrawPermissionClaimOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}
	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}

	//(fee)(withdraw_permission)(withdraw_from_account)(withdraw_to_account)(amount_to_withdraw)(memo)

	logging.Warnf("%s is not implemented", p.Type().OperationName())
	return nil
}
