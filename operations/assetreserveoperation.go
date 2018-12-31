package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeAssetReserve] = func() types.Operation {
		op := &AssetReserveOperation{}
		return op
	}
}

type AssetReserveOperation struct {
	types.OperationFee
	Payer           types.GrapheneID  `json:"payer"`
	AmountToReserve types.AssetAmount `json:"amount_to_reserve"`
	Extensions      types.Extensions  `json:"extensions"`
}

func (p AssetReserveOperation) Type() types.OperationType {
	return types.OperationTypeAssetReserve
}

func (p AssetReserveOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Payer); err != nil {
		return errors.Annotate(err, "encode payer")
	}

	if err := enc.Encode(p.AmountToReserve); err != nil {
		return errors.Annotate(err, "encode amount to reverse")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}
