package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeCallOrderUpdate] = func() types.Operation {
		op := &CallOrderUpdateOperation{}
		return op
	}
}

type CallOrderUpdateOperation struct {
	types.OperationFee
	DeltaCollateral types.AssetAmount               `json:"delta_collateral"`
	DeltaDebt       types.AssetAmount               `json:"delta_debt"`
	FundingAccount  types.GrapheneID                `json:"funding_account"`
	Extensions      types.CallOrderUpdateExtensions `json:"extensions"`
}

func (p CallOrderUpdateOperation) Type() types.OperationType {
	return types.OperationTypeCallOrderUpdate
}

func (p CallOrderUpdateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.FundingAccount); err != nil {
		return errors.Annotate(err, "encode funding account")
	}

	if err := enc.Encode(p.DeltaCollateral); err != nil {
		return errors.Annotate(err, "encode delta collateral")
	}

	if err := enc.Encode(p.DeltaDebt); err != nil {
		return errors.Annotate(err, "encode delta debt")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}
