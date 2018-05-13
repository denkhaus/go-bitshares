package operations

//go:generate ffjson   $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &AssetPublishFeedOperation{}
	types.OperationMap[op.Type()] = op
}

type AssetPublishFeedOperation struct {
	Publisher  types.GrapheneID  `json:"publisher"`
	AssetID    types.GrapheneID  `json:"asset_id"`
	Feed       types.PriceFeed   `json:"feed"`
	Fee        types.AssetAmount `json:"fee"`
	Extensions types.Extensions  `json:"extensions"`
}

func (p *AssetPublishFeedOperation) ApplyFee(fee types.AssetAmount) {
	p.Fee = fee
}

func (p AssetPublishFeedOperation) Type() types.OperationType {
	return types.OperationTypeAssetPublishFeed
}

func (p AssetPublishFeedOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Publisher); err != nil {
		return errors.Annotate(err, "encode publisher")
	}

	if err := enc.Encode(p.AssetID); err != nil {
		return errors.Annotate(err, "encode asset id")
	}

	if err := enc.Encode(p.Feed); err != nil {
		return errors.Annotate(err, "encode feed")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

//NewAssetPublishFeedOperation creates a new AssetPublishFeedOperation
func NewAssetPublishFeedOperation() *AssetPublishFeedOperation {
	tx := AssetPublishFeedOperation{
		Extensions: types.Extensions{},
	}
	return &tx
}
