package objects

//go:generate ffjson $GOFILE

type AssetFeedInfo struct {
	MaintenanceCollateralRatio UInt64 `json:"maintenance_collateral_ratio"`
	MaximumShortSqueezeRatio   UInt64 `json:"maximum_short_squeeze_ratio"`
	SettlementPrice            Price  `json:"settlement_price"`
	CoreExchangeRate           Price  `json:"core_exchange_rate"`
}
