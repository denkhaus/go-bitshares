package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &BalanceClaimOperation{}
	types.OperationMap[op.Type()] = op
}

type BalanceClaimOperation struct {
	BalanceToClaim   types.GrapheneID  `json:"balance_to_claim"`
	BalanceOwnerKey  types.PublicKey   `json:"balance_owner_key"`
	DepositToAccount types.GrapheneID  `json:"deposit_to_account"`
	TotalClaimed     types.AssetAmount `json:"total_claimed"`
	Fee              types.AssetAmount `json:"fee"`
}

func (p *BalanceClaimOperation) ApplyFee(fee types.AssetAmount) {
	p.Fee = fee
}

func (p BalanceClaimOperation) Type() types.OperationType {
	return types.OperationTypeBalanceClaim
}

func (p BalanceClaimOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.DepositToAccount); err != nil {
		return errors.Annotate(err, "encode DepositToAccount")
	}

	if err := enc.Encode(p.BalanceToClaim); err != nil {
		return errors.Annotate(err, "encode BalanceToClaim")
	}

	if err := enc.Encode(p.BalanceOwnerKey); err != nil {
		return errors.Annotate(err, "encode BalanceOwnerKey")
	}

	if err := enc.Encode(p.TotalClaimed); err != nil {
		return errors.Annotate(err, "encode TotalClaimed")
	}

	return nil
}

//NewBalanceClaimOperation creates a new BalanceClaimOperation
func NewBalanceClaimOperation() *BalanceClaimOperation {
	tx := BalanceClaimOperation{}
	return &tx
}
