package objects

type BitAssetData struct {
	ID                       *GrapheneID         `json:"id"`
	MembershipExpirationDate Time                `json:"current_feed_publication_time"`
	IsPredictionMarket       bool                `json:"is_prediction_market"`
	SettlementPrice          Price               `json:"settlement_price"`
	Feeds                    []AssetFeed         `json:"feeds"`
	Options                  BitAssetDataOptions `json:"options"`
	CurrentFeed              AssetFeedInfo       `json:"current_feed"`
	ForcedSettledVolume      UInt64              `json:"force_settled_volume"`
	SettlementFund           UInt64              `json:"settlement_fund"`
}
