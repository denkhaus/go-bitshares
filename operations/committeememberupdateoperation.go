package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeCommitteeMemberUpdate] = func() types.Operation {
		op := &CommitteeMemberUpdateOperation{}
		return op
	}
}

type CommitteeMemberUpdateOperation struct {
	types.OperationFee
	CommitteeMember        types.CommitteeMember `json:"committee_member"`
	CommitteeMemberAccount types.GrapheneID      `json:"committee_member_account"`
	NewURL                 *types.String         `json:"new_url,omitempty"`
}

func (p CommitteeMemberUpdateOperation) Type() types.OperationType {
	return types.OperationTypeCommitteeMemberUpdate
}

func (p CommitteeMemberUpdateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}
	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}
	if err := enc.Encode(p.CommitteeMember); err != nil {
		return errors.Annotate(err, "encode CommitteeMember")
	}
	if err := enc.Encode(p.CommitteeMemberAccount); err != nil {
		return errors.Annotate(err, "encode CommitteeMemberAccount")
	}
	if err := enc.Encode(p.NewURL != nil); err != nil {
		return errors.Annotate(err, "encode NewURL available")
	}
	if err := enc.Encode(p.NewURL); err != nil {
		return errors.Annotate(err, "encode NewURL")
	}

	return nil
}
