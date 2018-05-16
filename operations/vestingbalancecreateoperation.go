package operations

//go:generate ffjson   $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &VestingBalanceCreateOperation{}
	types.OperationMap[op.Type()] = op
}

type VestingBalanceCreatePolicy struct {
	StartClaim     types.Time   `json:"start_claim"`
	VestingSeconds types.UInt64 `json:"vesting_seconds"`
}

// Operation 32 ------------------------- dump start ---------------------------------------
// (map[string]interface {}) (len=5) {
//  (string) (len=6) "amount": (map[string]interface {}) (len=2) {
//   (string) (len=6) "amount": (float64) 1.73154449e+08,
//   (string) (len=8) "asset_id": (string) (len=8) "1.3.1564"
//  },
//  (string) (len=7) "creator": (string) (len=10) "1.2.913645",
//  (string) (len=3) "fee": (map[string]interface {}) (len=2) {
//   (string) (len=6) "amount": (float64) 1.736754e+07,
//   (string) (len=8) "asset_id": (string) (len=8) "1.3.1564"
//  },
//  (string) (len=5) "owner": (string) (len=10) "1.2.913645",
//  (string) (len=6) "policy": ([]interface {}) (len=2 cap=2) {
//   (float64) 1,
//   (map[string]interface {}) (len=2) {
//    (string) (len=11) "start_claim": (string) (len=19) "2018-05-10T15:07:55",
//    (string) (len=15) "vesting_seconds": (float64) 3.1536e+07
//   }
//  }
// }
// Operation 32 -------------------------  dump end  ---------------------------------------

//TODO: implement policy
type VestingBalanceCreateOperation struct {
	Amount  types.AssetAmount `json:"amount"`
	Creator types.GrapheneID  `json:"creator"`
	Fee     types.AssetAmount `json:"fee"`
	Owner   types.GrapheneID  `json:"owner"`
	Policy  []interface{}     `json:"policy"`
}

func (p *VestingBalanceCreateOperation) ApplyFee(fee types.AssetAmount) {
	p.Fee = fee
}

func (p VestingBalanceCreateOperation) Type() types.OperationType {
	return types.OperationTypeVestingBalanceCreate
}

//TODO: define!
func (p VestingBalanceCreateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	return nil
}

//NewVestingBalanceCreateOperation creates a new VestingBalanceCreateOperation
func NewVestingBalanceCreateOperation() *VestingBalanceCreateOperation {
	tx := VestingBalanceCreateOperation{}
	return &tx
}
