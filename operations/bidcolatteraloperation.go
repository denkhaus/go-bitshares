package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &BidColatteralOperation{}
	types.OperationMap[op.Type()] = op
}

type BidColatteralOperation struct {
	AdditionalCollateral types.AssetAmount `json:"additional_collateral"`
	Bidder               types.GrapheneID  `json:"bidder"`
	DebtCovered          types.AssetAmount `json:"debt_covered"`
	Extensions           types.Extensions  `json:"extensions"`
	Fee                  types.AssetAmount `json:"fee"`
}

func (p BidColatteralOperation) GetFee() types.AssetAmount {
	return p.Fee
}

func (p *BidColatteralOperation) SetFee(fee types.AssetAmount) {
	p.Fee = fee
}

func (p BidColatteralOperation) Type() types.OperationType {
	return types.OperationTypeBidColatteral
}

func (p BidColatteralOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Bidder); err != nil {
		return errors.Annotate(err, "encode Bidder")
	}

	if err := enc.Encode(p.AdditionalCollateral); err != nil {
		return errors.Annotate(err, "encode AdditionalCollateral")
	}

	if err := enc.Encode(p.DebtCovered); err != nil {
		return errors.Annotate(err, "encode DebtCovered")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extensions")
	}

	return nil
}

//NewBidColatteralOperation creates a new BidColatteralOperation
func NewBidColatteralOperation() *BidColatteralOperation {
	tx := BidColatteralOperation{
		Extensions: types.Extensions{},
	}
	return &tx
}
