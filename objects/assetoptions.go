package objects

const (
	CHARGE_MARKET_FEE    = 0x01
	WHITE_LIST           = 0x02
	OVERRIDE_AUTHORITY   = 0x04
	TRANSFER_RESTRICTED  = 0x08
	DISABLE_FORCE_SETTLE = 0x10
	GLOBAL_SETTLE        = 0x20
	DISABLE_CONFIDENTIAL = 0x40
	WITNESS_FED_ASSET    = 0x80
	COMITEE_FED_ASSET    = 0x100
)

//TODO: Implement whitelist_authorities, blacklist_authorities, whitelist_markets, blacklist_markets and extensions
type AssetOptions struct {
	MaxSupply         uint64  `json:"max_supply"`
	MarketFeePercent  float64 `json:"market_fee_percent"`
	MaxMarketFee      uint64  `json:"max_market_fee"`
	IssuerPermissions int64   `json:"issuer_permissions"`
	Flags             int     `json:"flags"`
	Description       string  `json:"description"`
	//coreExchangeRate  Price

}
