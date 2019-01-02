package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeAssetPublishFeed] = func() types.Operation {
		op := &AssetPublishFeedOperation{}
		return op
	}
}

type AssetPublishFeedOperation struct {
	types.OperationFee
	Publisher  types.AccountID  `json:"publisher"`
	AssetID    types.AssetID    `json:"asset_id"`
	Feed       types.PriceFeed  `json:"feed"`
	Extensions types.Extensions `json:"extensions"`
}

func (p AssetPublishFeedOperation) Type() types.OperationType {
	return types.OperationTypeAssetPublishFeed
}

func (p AssetPublishFeedOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}

	if err := enc.Encode(p.Publisher); err != nil {
		return errors.Annotate(err, "encode Publisher")
	}

	if err := enc.Encode(p.AssetID); err != nil {
		return errors.Annotate(err, "encode AssetID")
	}

	if err := enc.Encode(p.Feed); err != nil {
		return errors.Annotate(err, "encode Feed")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extensions")
	}

	return nil
}
