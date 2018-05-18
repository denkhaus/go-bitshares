//This file is generated by btsgen. DO NOT EDIT.
//operation sample data for OperationTypeAssetCreate

package samples

import (
	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/types"
)

func init() {
	data.OpSampleMap[types.OperationTypeAssetCreate] =
		sampleDataAssetCreateOperation
}

var sampleDataAssetCreateOperation = `{
  "bitasset_opts": {
    "extensions": [],
    "feed_lifetime_sec": 86400,
    "force_settlement_delay_sec": 86400,
    "force_settlement_offset_percent": 0,
    "maximum_force_settlement_volume": 2000,
    "minimum_feeds": 1,
    "short_backing_asset": "1.3.0"
  },
  "common_options": {
    "blacklist_authorities": [],
    "blacklist_markets": [],
    "core_exchange_rate": {
      "base": {
        "amount": 1,
        "asset_id": "1.3.1"
      },
      "quote": {
        "amount": 1,
        "asset_id": "1.3.0"
      }
    },
    "description": "",
    "extensions": [],
    "flags": 128,
    "issuer_permissions": 511,
    "market_fee_percent": 0,
    "max_market_fee": "1000000000000000",
    "max_supply": "1000000000000000",
    "whitelist_authorities": [],
    "whitelist_markets": []
  },
  "extensions": [],
  "fee": {
    "amount": "26000000001",
    "asset_id": "1.3.0"
  },
  "is_prediction_market": false,
  "issuer": "1.2.121",
  "precision": 4,
  "symbol": "TCNY"
}`

//end of file