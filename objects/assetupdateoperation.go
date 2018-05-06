package objects

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type AssetUpdateOperation struct {
	AssetToUpdate GrapheneID   `json:"asset_to_update"`
	Issuer        GrapheneID   `json:"issuer"`
	Fee           AssetAmount  `json:"fee"`
	Extensions    Extensions   `json:"extensions"`
	Options       AssetOptions `json:"new_options"`
}

//implements Operation interface
func (p *AssetUpdateOperation) ApplyFee(fee AssetAmount) {
	p.Fee = fee
}

//implements Operation interface
func (p AssetUpdateOperation) Type() OperationType {
	return OperationTypeAssetUpdate
}

//TODO: validate encode order!
//implements Operation interface
func (p AssetUpdateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.AssetToUpdate); err != nil {
		return errors.Annotate(err, "encode assettoupdate")
	}

	if err := enc.Encode(p.Issuer); err != nil {
		return errors.Annotate(err, "encode issuer")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	if err := enc.Encode(p.Options); err != nil {
		return errors.Annotate(err, "encode options")
	}

	return nil
}

//NewAssetUpdateOperation creates a new AssetUpdateOperation
func NewAssetUpdateOperation() *AssetUpdateOperation {
	tx := AssetUpdateOperation{
		Extensions: Extensions{},
	}
	return &tx
}
