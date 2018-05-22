package operations

//go:generate ffjson   $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &ProposalUpdateOperation{}
	types.OperationMap[op.Type()] = op
}

type ProposalUpdateOperation struct {
	ActiveApprovalsToAdd    types.GrapheneIDs `json:"active_approvals_to_add"`
	ActiveApprovalsToRemove types.GrapheneIDs `json:"active_approvals_to_remove"`
	Extensions              types.Extensions  `json:"extensions"`
	Fee                     types.AssetAmount `json:"fee"`
	FeePayingAccount        types.GrapheneID  `json:"fee_paying_account"`
	KeyApprovalsToAdd       types.GrapheneIDs `json:"key_approvals_to_add"`
	KeyApprovalsToRemove    types.GrapheneIDs `json:"key_approvals_to_remove"`
	OwnerApprovalsToAdd     types.GrapheneIDs `json:"owner_approvals_to_add"`
	OwnerApprovalsToRemove  types.GrapheneIDs `json:"owner_approvals_to_remove"`
	Proposal                types.GrapheneID  `json:"proposal"`
}

func (p *ProposalUpdateOperation) ApplyFee(fee types.AssetAmount) {
	p.Fee = fee
}

func (p ProposalUpdateOperation) Type() types.OperationType {
	return types.OperationTypeProposalUpdate
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

//NewProposalUpdateOperation creates a new ProposalUpdateOperation
func NewProposalUpdateOperation() *ProposalUpdateOperation {
	tx := ProposalUpdateOperation{
		Extensions: types.Extensions{},
	}
	return &tx
}
