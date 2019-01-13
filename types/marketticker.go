package types

type MarketTicker struct {
	Time          Time    `json:"time"`
	Base          AssetID `json:"base"`
	Quote         AssetID `json:"quote"`
	Latest        Float64 `json:"latest"`
	LowestAsk     Float64 `json:"lowest_ask"`
	HighestBid    Float64 `json:"highest_bid"`
	PercentChange Float64 `json:"percent_change"`
	BaseVolume    Float64 `json:"base_volume"`
	QuoteVolume   Float64 `json:"quote_volume"`
}
