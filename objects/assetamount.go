package objects

import (
	"math"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type AssetAmount struct {
	Asset  GrapheneID `json:"asset_id"`
	Amount Int64      `json:"amount"`
}

func (p *AssetAmount) Rate(prec float64) float64 {
	return float64(p.Amount) / math.Pow(10, prec)
}

func (p AssetAmount) Valid() bool {
	return p.Asset.Valid() && p.Amount != 0
}

//implements Operation interface
func (p AssetAmount) Marshal(enc *util.TypeEncoder) error {

	if err := enc.Encode(p.Amount); err != nil {
		return errors.Annotate(err, "encode amount")
	}

	if err := enc.Encode(p.Asset); err != nil {
		return errors.Annotate(err, "encode asset")
	}

	return nil
}
