//This file is generated by btsgen. DO NOT EDIT.
//operation sample data for OperationTypeCallOrderUpdate

package samples

import (
	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/types"
)

func init() {
	data.OpSampleMap[types.OperationTypeCallOrderUpdate] =
		sampleDataCallOrderUpdateOperation
}

var sampleDataCallOrderUpdateOperation = `{
  "delta_collateral": {
    "amount": "16868579331",
    "asset_id": "1.3.0"
  },
  "delta_debt": {
    "amount": 2500000,
    "asset_id": "1.3.121"
  },
  "extensions": {
    "target_collateral_ratio": 1750
  },
  "fee": {
    "amount": 4000000,
    "asset_id": "1.3.0"
  },
  "funding_account": "1.2.188"
}`
