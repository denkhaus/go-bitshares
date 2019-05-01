package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeProposalUpdate] = func() types.Operation {
		op := &ProposalUpdateOperation{}
		return op
	}
}

type ProposalUpdateOperation struct {
	types.OperationFee
	ActiveApprovalsToAdd    types.AccountIDs `json:"active_approvals_to_add"`
	ActiveApprovalsToRemove types.AccountIDs `json:"active_approvals_to_remove"`
	OwnerApprovalsToAdd     types.AccountIDs `json:"owner_approvals_to_add"`
	OwnerApprovalsToRemove  types.AccountIDs `json:"owner_approvals_to_remove"`
	Extensions              types.Extensions `json:"extensions"`
	FeePayingAccount        types.AccountID  `json:"fee_paying_account"`
	KeyApprovalsToAdd       types.PublicKeys `json:"key_approvals_to_add"`
	KeyApprovalsToRemove    types.PublicKeys `json:"key_approvals_to_remove"`
	Proposal                types.ProposalID `json:"proposal"`
}

func (p ProposalUpdateOperation) Type() types.OperationType {
	return types.OperationTypeProposalUpdate
}

func (p ProposalUpdateOperation) MarshalFeeScheduleParams(params types.M, enc *util.TypeEncoder) error {
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

func (p ProposalUpdateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}

	if err := enc.Encode(p.FeePayingAccount); err != nil {
		return errors.Annotate(err, "encode FeePayingAccount")
	}

	if err := enc.Encode(p.Proposal); err != nil {
		return errors.Annotate(err, "encode Proposal")
	}

	if err := enc.Encode(p.ActiveApprovalsToAdd); err != nil {
		return errors.Annotate(err, "encode ActiveApprovalsToAdd")
	}

	if err := enc.Encode(p.ActiveApprovalsToRemove); err != nil {
		return errors.Annotate(err, "encode ActiveApprovalsToRemove")
	}

	if err := enc.Encode(p.OwnerApprovalsToAdd); err != nil {
		return errors.Annotate(err, "encode OwnerApprovalsToAdd")
	}

	if err := enc.Encode(p.OwnerApprovalsToRemove); err != nil {
		return errors.Annotate(err, "encode OwnerApprovalsToRemove")
	}

	if err := enc.Encode(p.KeyApprovalsToAdd); err != nil {
		return errors.Annotate(err, "encode KeyApprovalsToAdd")
	}

	if err := enc.Encode(p.KeyApprovalsToRemove); err != nil {
		return errors.Annotate(err, "encode KeyApprovalsToRemove")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extensions")
	}

	return nil
}
