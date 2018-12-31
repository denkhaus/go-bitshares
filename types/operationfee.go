package types

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type OperationFee struct {
	Fee *AssetAmount `json:"fee,omitempty"`
}

func (p OperationFee) GetFee() AssetAmount {
	return *p.Fee
}

func (p *OperationFee) SetFee(fee AssetAmount) {
	p.Fee = &fee
}

func (p *OperationFee) MarshalFeeScheduleParams(params M, enc *util.TypeEncoder) error {
	if fee, ok := params["fee"]; ok {
		if err := enc.Encode(UInt64(fee.(float64))); err != nil {
			return errors.Annotate(err, "encode Fee")
		}
	}

	return nil
}
