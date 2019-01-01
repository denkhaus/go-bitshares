package types

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

//go:generate ffjson $GOFILE
type CommitteeMember struct {
	ID                     GrapheneID `json:"id"`
	CommitteeMemberAccount GrapheneID `json:"committee_member_account"`
	TotalVotes             UInt64     `json:"total_votes"`
	URL                    String     `json:"url"`
	VoteID                 VoteID     `json:"vote_id"`
}

func (p CommitteeMember) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.CommitteeMemberAccount); err != nil {
		return errors.Annotate(err, "encode CommitteeMemberAccount")
	}
	if err := enc.Encode(p.VoteID); err != nil {
		return errors.Annotate(err, "encode VoteID")
	}
	if err := enc.Encode(p.TotalVotes); err != nil {
		return errors.Annotate(err, "encode TotalVotes")
	}
	if err := enc.Encode(p.URL); err != nil {
		return errors.Annotate(err, "encode URL")
	}

	return nil
}
