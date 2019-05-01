package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeBidCollateral] = func() types.Operation {
		op := &BidCollateralOperation{}
		return op
	}
}

type BidCollateralOperation struct {
	types.OperationFee
	AdditionalCollateral types.AssetAmount `json:"additional_collateral"`
	Bidder               types.AccountID   `json:"bidder"`
	DebtCovered          types.AssetAmount `json:"debt_covered"`
	Extensions           types.Extensions  `json:"extensions"`
}

func (p BidCollateralOperation) Type() types.OperationType {
	return types.OperationTypeBidCollateral
}

func (p BidCollateralOperation) Marshal(enc *util.TypeEncoder) error {
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
