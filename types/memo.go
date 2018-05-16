package types

//go:generate ffjson   $GOFILE

import (
	"encoding/hex"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type Memo struct {
	From    PublicKey `json:"from"`
	To      PublicKey `json:"to"`
	Nonce   UInt64    `json:"nonce"`
	Message string    `json:"message"`
}

func (p Memo) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.From); err != nil {
		return errors.Annotate(err, "encode from")
	}

	if err := enc.Encode(p.To); err != nil {
		return errors.Annotate(err, "encode to")
	}

	if err := enc.Encode(p.Nonce); err != nil {
		return errors.Annotate(err, "encode nonce")
	}

	q, err := hex.DecodeString(p.Message)
	if err != nil {
		return errors.Annotate(err, "DecodeString")
	}

	if err := enc.Encode(string(q)); err != nil {
		return errors.Annotate(err, "encode message")
	}

	return nil
}
