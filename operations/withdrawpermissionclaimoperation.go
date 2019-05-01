package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
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
	WithdrawPermission  types.WithdrawPermissionID `json:"withdraw_permission"`
	WithdrawFromAccount types.AccountID            `json:"withdraw_from_account"`
	WithdrawToAccount   types.AccountID            `json:"withdraw_to_account"`
	AmountToWithdraw    types.AssetAmount          `json:"amount_to_withdraw"`
	Memo                *types.Memo                `json:"memo,omitempty"`
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
	if err := enc.Encode(p.WithdrawPermission); err != nil {
		return errors.Annotate(err, "encode WithdrawPermission")
	}
	if err := enc.Encode(p.WithdrawFromAccount); err != nil {
		return errors.Annotate(err, "encode WithdrawFromAccount")
	}
	if err := enc.Encode(p.WithdrawToAccount); err != nil {
		return errors.Annotate(err, "encode WithdrawToAccount")
	}
	if err := enc.Encode(p.AmountToWithdraw); err != nil {
		return errors.Annotate(err, "encode AmountToWithdraw")
	}
	if err := enc.Encode(p.Memo != nil); err != nil {
		return errors.Annotate(err, "encode Memo available")
	}
	if err := enc.Encode(p.Memo); err != nil {
		return errors.Annotate(err, "encode Memo")
	}

	return nil
}
