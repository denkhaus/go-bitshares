package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &AssetCreateOperation{}
	types.OperationMap[op.Type()] = op
}

type AssetCreateOperation struct {
	types.OperationFee
	BitassetOptions    *types.BitassetOptions `json:"bitasset_opts"`
	CommonOptions      types.AssetOptions     `json:"common_options"`
	Extensions         types.Extensions       `json:"extensions"`
	IsPredictionMarket bool                   `json:"is_prediction_market"`
	Issuer             types.GrapheneID       `json:"issuer"`
	Precision          types.UInt8            `json:"precision"`
	Symbol             string                 `json:"symbol"`
}

func (p AssetCreateOperation) Type() types.OperationType {
	return types.OperationTypeAssetCreate
}

func (p AssetCreateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Issuer); err != nil {
		return errors.Annotate(err, "encode issuer")
	}

	if err := enc.Encode(p.Symbol); err != nil {
		return errors.Annotate(err, "encode Symbol")
	}

	if err := enc.Encode(p.Precision); err != nil {
		return errors.Annotate(err, "encode Precision")
	}

	if err := enc.Encode(p.CommonOptions); err != nil {
		return errors.Annotate(err, "encode CommonOptions")
	}

	if err := enc.Encode(p.BitassetOptions != nil); err != nil {
		return errors.Annotate(err, "encode have BitassetOptions")
	}

	if err := enc.Encode(p.BitassetOptions); err != nil {
		return errors.Annotate(err, "encode BitassetOptions")
	}

	if err := enc.Encode(p.IsPredictionMarket); err != nil {
		return errors.Annotate(err, "encode IsPredictionMarket")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

//NewAssetCreateOperation creates a new AssetCreateOperation
func NewAssetCreateOperation() *AssetCreateOperation {
	tx := AssetCreateOperation{}
	return &tx
}
