package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeAssetClaimFees] = func() types.Operation {
		op := &AssetClaimFeesOperation{}
		return op
	}
}

type AssetClaimFeesOperation struct {
	types.OperationFee
	Issuer        types.AccountID   `json:"issuer"`
	AmountToClaim types.AssetAmount `json:"amount_to_claim"`
	Extensions    types.Extensions  `json:"extensions"`
}

func (p AssetClaimFeesOperation) Type() types.OperationType {
	return types.OperationTypeAssetClaimFees
}

func (p AssetClaimFeesOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation type")
	}
	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}
	if err := enc.Encode(p.Issuer); err != nil {
		return errors.Annotate(err, "encode Issuer")
	}
	if err := enc.Encode(p.AmountToClaim); err != nil {
		return errors.Annotate(err, "encode AmountToClaim")
	}
	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extensions")
	}
	return nil
}
