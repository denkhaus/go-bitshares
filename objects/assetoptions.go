package objects

// TODO: rename
const (
	AssetFlagCHARGE_MARKET_FEE    = 0x01
	AssetFlagWHITE_LIST           = 0x02
	AssetFlagOVERRIDE_AUTHORITY   = 0x04
	AssetFlagTRANSFER_RESTRICTED  = 0x08
	AssetFlagDISABLE_FORCE_SETTLE = 0x10
	AssetFlagGLOBAL_SETTLE        = 0x20
	AssetFlagDISABLE_CONFIDENTIAL = 0x40
	AssetFlagWITNESS_FED_ASSET    = 0x80
	AssetFlagCOMITEE_FED_ASSET    = 0x100
)

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
	Extensions           []interface{} `json:"extensions"`
}
