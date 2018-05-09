// (map[string]interface {}) (len=5) {
// 	(string) (len=17) "balance_owner_key": (string) (len=53) "BTS5ZZfxSKoPck2kGmZzfzmCufHE48B2NAV71iFa8UXj2BPkaVr7P",
// 	(string) (len=16) "balance_to_claim": (string) (len=10) "1.15.69075",
// 	(string) (len=18) "deposit_to_account": (string) (len=9) "1.2.10408",
// 	(string) (len=3) "fee": (map[string]interface {}) (len=2) {
// 	 (string) (len=6) "amount": (float64) 0,
// 	 (string) (len=8) "asset_id": (string) (len=5) "1.3.0"
// 	},
// 	(string) (len=13) "total_claimed": (map[string]interface {}) (len=2) {
// 	 (string) (len=6) "amount": (float64) 1.22768008e+08,
// 	 (string) (len=8) "asset_id": (string) (len=5) "1.3.0"
// 	}
//    }

package objects

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &BalanceClaimOperation{}
	opMap[op.Type()] = op
}

type BalanceClaimOperation struct {
	BalanceToClaim   GrapheneID  `json:"balance_to_claim"`
	BalanceOwnerKey  PublicKey   `json:"balance_owner_key"`
	DepositToAccount GrapheneID  `json:"deposit_to_account"`
	TotalClaimed     AssetAmount `json:"total_claimed"`
	Fee              AssetAmount `json:"fee"`
}

//implements Operation interface
func (p *BalanceClaimOperation) ApplyFee(fee AssetAmount) {
	p.Fee = fee
}

//implements Operation interface
func (p BalanceClaimOperation) Type() OperationType {
	return OperationTypeBalanceClaim
}

//TODO: validate encode order!
//implements Operation interface
func (p BalanceClaimOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.BalanceOwnerKey); err != nil {
		return errors.Annotate(err, "encode balance owner key")
	}

	if err := enc.Encode(p.BalanceToClaim); err != nil {
		return errors.Annotate(err, "encode balance to claim")
	}

	if err := enc.Encode(p.DepositToAccount); err != nil {
		return errors.Annotate(err, "encode deposit to account")
	}

	if err := enc.Encode(p.TotalClaimed); err != nil {
		return errors.Annotate(err, "encode total claimed")
	}

	return nil
}

//NewBalanceClaimOperation creates a new BalanceClaimOperation
func NewBalanceClaimOperation() *BalanceClaimOperation {
	tx := BalanceClaimOperation{}
	return &tx
}
