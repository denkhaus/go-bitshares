package objects

type BitAssetDataOptions struct {
	ShortBackingAsset            GrapheneID `json:"short_backing_asset"`
	MinimumFeeds                 UInt64     `json:"minimum_feeds"`
	ForceSettlementDelaySec      UInt64     `json:"force_settlement_delay_sec"`
	ForceSettlementOffsetPercent UInt64     `json:"force_settlement_offset_percent"`
	MaximumForceSettlementVolume UInt64     `json:"maximum_force_settlement_volume"`
	Extensions                   Extensions `json:"extensions"`
	FeedLifetimeSec              UInt64     `json:"feed_lifetime_sec"`
}
