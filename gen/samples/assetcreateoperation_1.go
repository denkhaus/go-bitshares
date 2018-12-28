//This file is generated by btsgen. DO NOT EDIT.
//operation sample data for OperationTypeAssetCreate

package samples

func init() {

	sampleDataAssetCreateOperation[1] = `{
  "bitasset_opts": {
    "extensions": [],
    "feed_lifetime_sec": 86400,
    "force_settlement_delay_sec": 86400,
    "force_settlement_offset_percent": 100,
    "maximum_force_settlement_volume": 2000,
    "minimum_feeds": 7,
    "short_backing_asset": "1.3.0"
  },
  "common_options": {
    "blacklist_authorities": [],
    "blacklist_markets": [],
    "core_exchange_rate": {
      "base": {
        "amount": 156000,
        "asset_id": "1.3.0"
      },
      "quote": {
        "amount": 100000,
        "asset_id": "1.3.1"
      }
    },
    "description": "{\"main\":\"The E-Krona is an electronic token which aims to track the value of the Swedish Krona so that 1 E-Krona equals the value of 1 Swedish Krona.\",\"short_name\":\"E-Krona\",\"market\":\"BTS\"}",
    "extensions": [],
    "flags": 0,
    "issuer_permissions": 507,
    "market_fee_percent": 0,
    "max_market_fee": 0,
    "max_supply": "10000000000",
    "whitelist_authorities": [],
    "whitelist_markets": []
  },
  "extensions": [],
  "fee": {
    "amount": 60701624,
    "asset_id": "1.3.0"
  },
  "is_prediction_market": false,
  "issuer": "1.2.171843",
  "precision": 5,
  "symbol": "EKRONA"
}`

}

//end of file
