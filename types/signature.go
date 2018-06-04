package types

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

//go:generate ffjson $GOFILE

type Signatures []Buffer

func (p Signatures) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	for _, sig := range p {
		if err := enc.Encode(sig.Bytes()); err != nil {
			return errors.Annotate(err, "encode Signature")
		}
	}

	return nil
}

func (p *Signatures) Reset() {
	*p = []Buffer{}
}
