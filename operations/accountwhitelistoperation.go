package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeAccountWhitelist] = func() types.Operation {
		op := &AccountWhitelistOperation{}
		return op
	}
}

type AccountWhitelistOperation struct {
	types.OperationFee
	AccountToList      types.GrapheneID `json:"account_to_list"`
	AuthorizingAccount types.GrapheneID `json:"authorizing_account"`
	Extensions         types.Extensions `json:"extensions"`
	NewListing         types.UInt8      `json:"new_listing"`
}

func (p AccountWhitelistOperation) Type() types.OperationType {
	return types.OperationTypeAccountWhitelist
}

//TODO: validate order
func (p AccountWhitelistOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}

	if err := enc.Encode(p.AuthorizingAccount); err != nil {
		return errors.Annotate(err, "encode AuthorizingAccount")
	}

	if err := enc.Encode(p.AccountToList); err != nil {
		return errors.Annotate(err, "encode AccountToList")
	}

	if err := enc.Encode(p.NewListing); err != nil {
		return errors.Annotate(err, "encode NewListing")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}
