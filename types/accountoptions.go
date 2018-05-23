package types

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type AccountOptions struct {
	MemoKey       PublicKey  `json:"memo_key"`
	VotingAccount GrapheneID `json:"voting_account"`
	NumWitness    UInt16     `json:"num_witness"`
	NumCommittee  UInt16     `json:"num_committee"`
	Votes         Votes      `json:"votes"`
	Extensions    Extensions `json:"extensions"`
}

func (p AccountOptions) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.MemoKey); err != nil {
		return errors.Annotate(err, "encode MemoKey")
	}

	if err := enc.Encode(p.VotingAccount); err != nil {
		return errors.Annotate(err, "encode VotingAccount")
	}

	if err := enc.Encode(p.NumWitness); err != nil {
		return errors.Annotate(err, "encode NumWitness")
	}

	if err := enc.Encode(p.NumCommittee); err != nil {
		return errors.Annotate(err, "encode NumCommittee")
	}

	if err := enc.Encode(p.Votes); err != nil {
		return errors.Annotate(err, "encode Votes")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extensions")
	}

	return nil
}
