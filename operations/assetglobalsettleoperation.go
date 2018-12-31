package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeAssetGlobalSettle] = func() types.Operation {
		op := &AssetGlobalSettleOperation{}
		return op
	}
}

type AssetGlobalSettleOperation struct {
	types.OperationFee
	AssetToSettle types.GrapheneID `json:"asset_to_settle"`
	Extensions    types.Extensions `json:"extensions"`
	Issuer        types.GrapheneID `json:"issuer"`
	SettlePrice   types.Price      `json:"settle_price"`
}

func (p AssetGlobalSettleOperation) Type() types.OperationType {
	return types.OperationTypeAssetGlobalSettle
}

func (p AssetGlobalSettleOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}
	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}
	if err := enc.Encode(p.Issuer); err != nil {
		return errors.Annotate(err, "encode Issuer")
	}
	if err := enc.Encode(p.AssetToSettle); err != nil {
		return errors.Annotate(err, "encode AssetToSettle")
	}
	if err := enc.Encode(p.SettlePrice); err != nil {
		return errors.Annotate(err, "encode SettlePrice")
	}
	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extensions")
	}
	return nil
}
