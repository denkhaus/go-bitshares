package operations

//go:generate ffjson  $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &FillOrderOperation{}
	types.OperationMap[op.Type()] = op
}

type FillOrderOperation struct {
	OrderID   types.GrapheneID  `json:"order_id"`
	AccountID types.GrapheneID  `json:"account_id"`
	Pays      types.AssetAmount `json:"pays"`
	Receives  types.AssetAmount `json:"receives"`
	Fee       types.AssetAmount `json:"fee"`
	IsMaker   bool              `json:"is_maker"`
	FillPrice types.Price       `json:"fill_price"`
}

func (p *FillOrderOperation) ApplyFee(fee types.AssetAmount) {
	p.Fee = fee
}

func (p FillOrderOperation) Type() types.OperationType {
	return types.OperationTypeFillOrder
}

//TODO: validate encode order!

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
