package types

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

//go:generate ffjson   $GOFILE

type PriceFeed struct {
	MaintenanceCollateralRatio UInt16 `json:"maintenance_collateral_ratio"`
	MaximumShortSqueezeRatio   UInt16 `json:"maximum_short_squeeze_ratio"`
	SettlementPrice            Price  `json:"settlement_price"`
	CoreExchangeRate           Price  `json:"core_exchange_rate"`
}

func (p PriceFeed) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.SettlementPrice); err != nil {
		return errors.Annotate(err, "encode SettlementPrice")
	}

	if err := enc.Encode(p.MaintenanceCollateralRatio); err != nil {
		return errors.Annotate(err, "encode MaintenanceCollateralRatio")
	}

	if err := enc.Encode(p.MaximumShortSqueezeRatio); err != nil {
		return errors.Annotate(err, "encode MaximumShortSqueezeRatio")
	}

	if err := enc.Encode(p.CoreExchangeRate); err != nil {
		return errors.Annotate(err, "encode CoreExchangeRate")
	}
	return nil
}
