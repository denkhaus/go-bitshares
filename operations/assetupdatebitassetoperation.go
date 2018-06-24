package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &AssetUpdateBitassetOperation{}
	types.OperationMap[op.Type()] = op
}

type AssetUpdateBitassetOperation struct {
	AssetToUpdate types.GrapheneID      `json:"asset_to_update"`
	Issuer        types.GrapheneID      `json:"issuer"`
	Fee           types.AssetAmount     `json:"fee"`
	Extensions    types.Extensions      `json:"extensions"`
	NewOptions    types.BitassetOptions `json:"new_options"`
}

func (p AssetUpdateBitassetOperation) GetFee() types.AssetAmount {
	return p.Fee
}

func (p *AssetUpdateBitassetOperation) SetFee(fee types.AssetAmount) {
	p.Fee = fee
}

func (p AssetUpdateBitassetOperation) Type() types.OperationType {
	return types.OperationTypeAssetUpdateBitasset
}

//TODO: validate order
func (p AssetUpdateBitassetOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Issuer); err != nil {
		return errors.Annotate(err, "encode issuer")
	}

	if err := enc.Encode(p.AssetToUpdate); err != nil {
		return errors.Annotate(err, "encode assettoupdate")
	}

	if err := enc.Encode(p.NewOptions); err != nil {
		return errors.Annotate(err, "encode new options")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

//NewAssetUpdateBitassetOperation creates a new AssetUpdateBitassetOperation
func NewAssetUpdateBitassetOperation() *AssetUpdateBitassetOperation {
	tx := AssetUpdateBitassetOperation{
		Extensions: types.Extensions{},
	}
	return &tx
}
