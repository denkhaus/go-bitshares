package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeAccountCreate] = func() types.Operation {
		op := &AccountCreateOperation{}
		return op
	}
}

type AccountCreateOperation struct {
	types.OperationFee
	Registrar       types.AccountID               `json:"registrar"`
	Referrer        types.AccountID               `json:"referrer"`
	ReferrerPercent types.UInt16                  `json:"referrer_percent"`
	Owner           types.Authority               `json:"owner"`
	Active          types.Authority               `json:"active"`
	Name            types.String                  `json:"name"`
	Extensions      types.AccountCreateExtensions `json:"extensions"`
	Options         types.AccountOptions          `json:"options"`
}

func (p AccountCreateOperation) Type() types.OperationType {
	return types.OperationTypeAccountCreate
}

func (p AccountCreateOperation) MarshalFeeScheduleParams(params types.M, enc *util.TypeEncoder) error {
	if bfee, ok := params["basic_fee"]; ok {
		if err := enc.Encode(types.UInt64(bfee.(float64))); err != nil {
			return errors.Annotate(err, "encode BasicFee")
		}
	}
	if pfee, ok := params["premium_fee"]; ok {
		if err := enc.Encode(types.UInt64(pfee.(float64))); err != nil {
			return errors.Annotate(err, "encode PremiumFee")
		}
	}
	if ppk, ok := params["price_per_kbyte"]; ok {
		if err := enc.Encode(types.UInt32(ppk.(float64))); err != nil {
			return errors.Annotate(err, "encode PricePerKByte")
		}
	}

	return nil
}

func (p AccountCreateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode Type")
	}
	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}
	if err := enc.Encode(p.Registrar); err != nil {
		return errors.Annotate(err, "encode Registrar")
	}
	if err := enc.Encode(p.Referrer); err != nil {
		return errors.Annotate(err, "encode Referrer")
	}
	if err := enc.Encode(p.ReferrerPercent); err != nil {
		return errors.Annotate(err, "encode ReferrerPercent")
	}
	if err := enc.Encode(p.Name); err != nil {
		return errors.Annotate(err, "encode Name")
	}
	if err := enc.Encode(p.Owner); err != nil {
		return errors.Annotate(err, "encode Owner")
	}
	if err := enc.Encode(p.Active); err != nil {
		return errors.Annotate(err, "encode Active")
	}
	if err := enc.Encode(p.Options); err != nil {
		return errors.Annotate(err, "encode Options")
	}
	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extensions")
	}

	return nil
}
