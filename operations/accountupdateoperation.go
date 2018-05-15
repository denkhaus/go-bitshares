package operations

//go:generate ffjson   $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	op := &AccountUpdateOperation{}
	types.OperationMap[op.Type()] = op
}

type AccountUpdateOperation struct {
	Account    types.GrapheneID     `json:"account"`
	Active     types.Authority      `json:"active"`
	Extensions types.Extensions     `json:"extensions"`
	Fee        types.AssetAmount    `json:"fee"`
	NewOptions types.AccountOptions `json:"new_options"`
	Owner      types.Authority      `json:"owner"`
}

func (p *AccountUpdateOperation) ApplyFee(fee types.AssetAmount) {
	p.Fee = fee
}

func (p AccountUpdateOperation) Type() types.OperationType {
	return types.OperationTypeAccountUpdate
}

//TODO: validate order
func (p AccountUpdateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation id")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}

	if err := enc.Encode(p.Account); err != nil {
		return errors.Annotate(err, "encode Account")
	}

	if err := enc.Encode(p.Active); err != nil {
		return errors.Annotate(err, "encode Active")
	}

	if err := enc.Encode(p.Owner); err != nil {
		return errors.Annotate(err, "encode Active")
	}

	if err := enc.Encode(p.NewOptions); err != nil {
		return errors.Annotate(err, "encode NewOptions")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

//NewAccountUpdateOperation creates a new AccountUpdateOperation
func NewAccountUpdateOperation() *AccountUpdateOperation {
	tx := AccountUpdateOperation{
		Extensions: types.Extensions{},
	}
	return &tx
}
