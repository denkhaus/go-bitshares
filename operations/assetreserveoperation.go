package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &AssetReserveOperation{}
	types.OperationMap[op.Type()] = op
}

type AssetReserveOperation struct {
	Payer           types.GrapheneID  `json:"payer"`
	AmountToReserve types.AssetAmount `json:"amount_to_reserve"`
	Fee             types.AssetAmount `json:"fee"`
	Extensions      types.Extensions  `json:"extensions"`
}

func (p AssetReserveOperation) GetFee() types.AssetAmount {
	return p.Fee
}

func (p *AssetReserveOperation) SetFee(fee types.AssetAmount) {
	p.Fee = fee
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

//NewAssetReserveOperation creates a new AssetReserveOperation
func NewAssetReserveOperation() *AssetReserveOperation {
	tx := AssetReserveOperation{
		Extensions: types.Extensions{},
	}
	return &tx
}
