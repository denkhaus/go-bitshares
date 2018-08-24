package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &ProposalDeleteOperation{}
	types.OperationMap[op.Type()] = op
}

type ProposalDeleteOperation struct {
	types.OperationFee
	Extensions          types.Extensions `json:"extensions"`
	FeePayingAccount    types.GrapheneID `json:"fee_paying_account"`
	Proposal            types.GrapheneID `json:"proposal"`
	UsingOwnerAuthority bool             `json:"using_owner_authority"`
}

func (p ProposalDeleteOperation) Type() types.OperationType {
	return types.OperationTypeProposalDelete
}

func (p ProposalDeleteOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}

	if err := enc.Encode(p.FeePayingAccount); err != nil {
		return errors.Annotate(err, "encode FeePayingAccount")
	}

	if err := enc.Encode(p.UsingOwnerAuthority); err != nil {
		return errors.Annotate(err, "encode UsingOwnerAuthority")
	}

	if err := enc.Encode(p.Proposal); err != nil {
		return errors.Annotate(err, "encode Proposal")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extensions")
	}

	return nil
}

//NewProposalDeleteOperation creates a new ProposalDeleteOperation
func NewProposalDeleteOperation() *ProposalDeleteOperation {
	tx := ProposalDeleteOperation{
		Extensions: types.Extensions{},
	}
	return &tx
}
