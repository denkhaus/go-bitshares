package types

//go:generate ffjson   $GOFILE

type AssetOptions struct {
	MaxSupply            UInt64        `json:"max_supply"`
	MaxMarketFee         UInt64        `json:"max_market_fee"`
	MarketFeePercent     UInt64        `json:"market_fee_percent"`
	Flags                int           `json:"flags"`
	Description          string        `json:"description"`
	CoreExchangeRate     Price         `json:"core_exchange_rate"`
	IssuerPermissions    int64         `json:"issuer_permissions"`
	BlacklistAuthorities []interface{} `json:"blacklist_authorities"`
	WhitelistAuthorities []interface{} `json:"whitelist_authorities"`
	BlacklistMarkets     []interface{} `json:"blacklist_markets"`
	WhitelistMarkets     []interface{} `json:"whitelist_markets"`
	Extensions           Extensions    `json:"extensions"`
}
