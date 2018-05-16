package types

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

//go:generate ffjson   $GOFILE

type BitassetOptions struct {
	FeedLifetimeSec              UInt32     `json:"feed_lifetime_sec"`
	MinimumFeeds                 UInt8      `json:"minimum_feeds"`
	ForceSettlementDelaySec      UInt32     `json:"force_settlement_delay_sec"`
	ForceSettlementOffsetPercent UInt16     `json:"force_settlement_offset_percent"`
	MaximumForceSettlementVolume UInt16     `json:"maximum_force_settlement_volume"`
	ShortBackingAsset            GrapheneID `json:"short_backing_asset"`
	Extensions                   Extensions `json:"extensions"`
}

func (p BitassetOptions) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.FeedLifetimeSec); err != nil {
		return errors.Annotate(err, "encode FeedLifetimeSec")
	}

	if err := enc.Encode(p.MinimumFeeds); err != nil {
		return errors.Annotate(err, "encode MinimumFeeds")
	}

	if err := enc.Encode(p.ForceSettlementDelaySec); err != nil {
		return errors.Annotate(err, "encode ForceSettlementDelaySec")
	}

	if err := enc.Encode(p.ForceSettlementOffsetPercent); err != nil {
		return errors.Annotate(err, "encode ForceSettlementOffsetPercent")
	}

	if err := enc.Encode(p.MaximumForceSettlementVolume); err != nil {
		return errors.Annotate(err, "encode MaximumForceSettlementVolume")
	}

	if err := enc.Encode(p.ShortBackingAsset); err != nil {
		return errors.Annotate(err, "encode ShortBackingAsset")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extensions")
	}

	return nil
}

type BitAssetData struct {
	ID                       GrapheneID      `json:"id"`
	MembershipExpirationDate Time            `json:"current_feed_publication_time"`
	IsPredictionMarket       bool            `json:"is_prediction_market"`
	SettlementPrice          Price           `json:"settlement_price"`
	Feeds                    AssetFeeds      `json:"feeds"`
	Options                  BitassetOptions `json:"options"`
	CurrentFeed              PriceFeed       `json:"current_feed"`
	ForcedSettledVolume      UInt64          `json:"force_settled_volume"`
	SettlementFund           UInt64          `json:"settlement_fund"`
}
