package objects

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

//go:generate ffjson -force-regenerate $GOFILE

func init() {
	op := &FillOrderOperation{}
	opMap[op.Type()] = op
}

type FillOrderOperation struct {
	OrderID   GrapheneID  `json:"order_id"`
	AccountID GrapheneID  `json:"account_id"`
	Pays      AssetAmount `json:"pays"`
	Receives  AssetAmount `json:"receives"`
	Fee       AssetAmount `json:"fee"`
	IsMaker   bool        `json:"is_maker"`
	FillPrice Price       `json:"fill_price"`
}

//implements Operation interface
func (p *FillOrderOperation) ApplyFee(fee AssetAmount) {
	p.Fee = fee
}

//implements Operation interface
func (p FillOrderOperation) Type() OperationType {
	return OperationTypeFillOrder
}

//TODO: validate encode order!
//implements Operation interface
func (p FillOrderOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.OrderID); err != nil {
		return errors.Annotate(err, "encode orderid")
	}

	if err := enc.Encode(p.AccountID); err != nil {
		return errors.Annotate(err, "encode accountid")
	}

	if err := enc.Encode(p.Pays); err != nil {
		return errors.Annotate(err, "encode pays")
	}

	if err := enc.Encode(p.Receives); err != nil {
		return errors.Annotate(err, "encode receives")
	}

	if err := enc.Encode(p.IsMaker); err != nil {
		return errors.Annotate(err, "encode ismaker")
	}

	if err := enc.Encode(p.FillPrice); err != nil {
		return errors.Annotate(err, "encode fillprice")
	}

	return nil
}

//NewFillOrderOperation creates a new FillOrderOperation
func NewFillOrderOperation() *FillOrderOperation {
	return &FillOrderOperation{}
}
