package types

//go:generate ffjson   $GOFILE

import (
	"math"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type AssetAmounts []AssetAmount

type AssetAmount struct {
	Amount Int64      `json:"amount"`
	Asset  GrapheneID `json:"asset_id"`
}

func (p *AssetAmount) Rate(prec float64) float64 {
	return float64(p.Amount) / math.Pow(10, prec)
}

func (p AssetAmount) Valid() bool {
	return p.Asset.Valid() && p.Amount != 0
}

func (p AssetAmount) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.Amount); err != nil {
		return errors.Annotate(err, "encode amount")
	}

	if err := enc.Encode(p.Asset); err != nil {
		return errors.Annotate(err, "encode asset")
	}

	return nil
}
