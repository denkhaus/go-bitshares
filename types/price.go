package types

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type Price struct {
	Base  AssetAmount `json:"base"`
	Quote AssetAmount `json:"quote"`
}

func (p Price) Rate(precBase, precQuote float64) Rate {
	return Rate(p.Base.Rate(precBase) / p.Quote.Rate(precQuote))
}

func (p Price) Valid() bool {
	return p.Base.Valid() && p.Quote.Valid()
}

func (p Price) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.Base); err != nil {
		return errors.Annotate(err, "encode amount")
	}

	if err := enc.Encode(p.Quote); err != nil {
		return errors.Annotate(err, "encode asset")
	}

	return nil
}
