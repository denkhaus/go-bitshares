package types

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

//go:generate ffjson $GOFILE

type AssetOptions struct {
	MaxSupply            Int64       `json:"max_supply"`
	MaxMarketFee         Int64       `json:"max_market_fee"`
	MarketFeePercent     UInt16      `json:"market_fee_percent"`
	Flags                UInt16      `json:"flags"`
	Description          string      `json:"description"`
	CoreExchangeRate     Price       `json:"core_exchange_rate"`
	IssuerPermissions    UInt16      `json:"issuer_permissions"`
	BlacklistAuthorities GrapheneIDs `json:"blacklist_authorities"`
	WhitelistAuthorities GrapheneIDs `json:"whitelist_authorities"`
	BlacklistMarkets     GrapheneIDs `json:"blacklist_markets"`
	WhitelistMarkets     GrapheneIDs `json:"whitelist_markets"`
	Extensions           Extensions  `json:"extensions"`
}

func (p AssetOptions) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.MaxSupply); err != nil {
		return errors.Annotate(err, "encode MaxSupply")
	}

	if err := enc.Encode(p.MarketFeePercent); err != nil {
		return errors.Annotate(err, "encode MarketFeePercent")
	}

	if err := enc.Encode(p.MaxMarketFee); err != nil {
		return errors.Annotate(err, "encode MaxMarketFee")
	}

	if err := enc.Encode(p.IssuerPermissions); err != nil {
		return errors.Annotate(err, "encode IssuerPermissions")
	}

	if err := enc.Encode(p.Flags); err != nil {
		return errors.Annotate(err, "encode Flags")
	}

	if err := enc.Encode(p.CoreExchangeRate); err != nil {
		return errors.Annotate(err, "encode CoreExchangeRate")
	}

	if err := enc.Encode(p.WhitelistAuthorities); err != nil {
		return errors.Annotate(err, "encode WhitelistAuthorities")
	}

	if err := enc.Encode(p.BlacklistAuthorities); err != nil {
		return errors.Annotate(err, "encode BlacklistAuthorities")
	}

	if err := enc.Encode(p.WhitelistMarkets); err != nil {
		return errors.Annotate(err, "encode WhitelistMarkets")
	}

	if err := enc.Encode(p.BlacklistMarkets); err != nil {
		return errors.Annotate(err, "encode BlacklistMarkets")
	}

	if err := enc.Encode(p.Description); err != nil {
		return errors.Annotate(err, "encode Description")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}
