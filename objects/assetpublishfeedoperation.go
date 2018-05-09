package objects

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &AssetPublishFeedOperation{}
	opMap[op.Type()] = op
}

type AssetPublishFeedOperation struct {
	Publisher GrapheneID    `json:"publisher"`
	AssetID   GrapheneID    `json:"asset_id"`
	Feed      AssetFeedInfo `json:"feed"`
	Fee       AssetAmount   `json:"fee"`

	Extensions Extensions `json:"extensions"`
}

//implements Operation interface
func (p *AssetPublishFeedOperation) ApplyFee(fee AssetAmount) {
	p.Fee = fee
}

//implements Operation interface
func (p AssetPublishFeedOperation) Type() OperationType {
	return OperationTypeAssetPublishFeed
}

//TODO: validate encode order!
//implements Operation interface
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
		Extensions: Extensions{},
	}
	return &tx
}
