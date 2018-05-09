package objects

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &CallOrderUpdateOperation{}
	opMap[op.Type()] = op
}

type CallOrderUpdateOperation struct {
	DeltaCollateral AssetAmount `json:"delta_collateral"`
	DeltaDebt       AssetAmount `json:"delta_debt"`
	FundingAccount  GrapheneID  `json:"funding_account"`
	Fee             AssetAmount `json:"fee"`
	Extensions      Extensions  `json:"extensions"`
}

//implements Operation interface
func (p *CallOrderUpdateOperation) ApplyFee(fee AssetAmount) {
	p.Fee = fee
}

//implements Operation interface
func (p CallOrderUpdateOperation) Type() OperationType {
	return OperationTypeCallOrderUpdate
}

//implements Operation interface
//order checked!
func (p CallOrderUpdateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
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

func NewCallOrderUpdateOperation(acct GrapheneID) *CallOrderUpdateOperation {
	op := CallOrderUpdateOperation{
		Extensions:     Extensions{},
		FundingAccount: acct,
	}

	return &op
}
