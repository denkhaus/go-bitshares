package types

//go:generate ffjson $GOFILE

import (
	"time"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type Transaction struct {
	RefBlockNum    UInt16     `json:"ref_block_num"`
	RefBlockPrefix UInt32     `json:"ref_block_prefix"`
	Expiration     Time       `json:"expiration"`
	Operations     Operations `json:"operations"`
	Extensions     Extensions `json:"extensions"`
}

func (p Transaction) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.RefBlockNum); err != nil {
		return errors.Annotate(err, "encode RefBlockNum")
	}

	if err := enc.Encode(p.RefBlockPrefix); err != nil {
		return errors.Annotate(err, "encode RefBlockPrefix")
	}

	if err := enc.Encode(p.Expiration); err != nil {
		return errors.Annotate(err, "encode Expiration")
	}

	if err := enc.Encode(p.Operations); err != nil {
		return errors.Annotate(err, "encode Operations")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extension")
	}

	return nil
}

//AdjustExpiration extends expiration by given duration.
func (p *Transaction) AdjustExpiration(dur time.Duration) {
	p.Expiration = p.Expiration.Add(dur)
}
