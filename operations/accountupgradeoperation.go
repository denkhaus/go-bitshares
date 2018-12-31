package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeAccountUpgrade] = func() types.Operation {
		op := &AccountUpgradeOperation{}
		return op
	}
}

type AccountUpgradeOperation struct {
	types.OperationFee
	AccountToUpgrade        types.GrapheneID `json:"account_to_upgrade"`
	Extensions              types.Extensions `json:"extensions"`
	UpgradeToLifetimeMember bool             `json:"upgrade_to_lifetime_member"`
}

func (p AccountUpgradeOperation) Type() types.OperationType {
	return types.OperationTypeAccountUpgrade
}

func (p AccountUpgradeOperation) MarshalFeeScheduleParams(params types.M, enc *util.TypeEncoder) error {
	if maf, ok := params["membership_annual_fee"]; ok {
		if err := enc.Encode(types.UInt64(maf.(float64))); err != nil {
			return errors.Annotate(err, "encode MembershipAnnualFee")
		}
	}
	if mlf, ok := params["membership_lifetime_fee"]; ok {
		if err := enc.Encode(types.UInt64(mlf.(float64))); err != nil {
			return errors.Annotate(err, "encode MembershipLifetimeFee")
		}
	}

	return nil
}

func (p AccountUpgradeOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}

	if err := enc.Encode(p.AccountToUpgrade); err != nil {
		return errors.Annotate(err, "encode AccountToUpgrade")
	}

	if err := enc.Encode(p.UpgradeToLifetimeMember); err != nil {
		return errors.Annotate(err, "encode UpgradeToLifetimeMember")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}
