package objects

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &AssetReverseOperation{}
	opMap[op.Type()] = op
}

type AssetReverseOperation struct {
	Payer           GrapheneID  `json:"payer"`
	AmountToReserve AssetAmount `json:"amount_to_reserve"`
	Fee             AssetAmount `json:"fee"`
	Extensions      Extensions  `json:"extensions"`
}

//implements Operation interface
func (p *AssetReverseOperation) ApplyFee(fee AssetAmount) {
	p.Fee = fee
}

//implements Operation interface
func (p AssetReverseOperation) Type() OperationType {
	return OperationTypeAssetReverse
}

//TODO: validate encode order!
//implements Operation interface
func (p AssetReverseOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Payer); err != nil {
		return errors.Annotate(err, "encode payer")
	}

	if err := enc.Encode(p.AmountToReserve); err != nil {
		return errors.Annotate(err, "encode amount to reserve")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

//NewAssetReverseOperation creates a new AssetReverseOperation
func NewAssetReverseOperation() *AssetReverseOperation {
	tx := AssetReverseOperation{
		Extensions: Extensions{},
	}
	return &tx
}
