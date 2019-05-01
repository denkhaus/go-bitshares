package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeTransferToBlind] = func() types.Operation {
		op := &TransferToBlindOperation{}
		return op
	}
}

type TransferToBlindOperation struct {
	types.OperationFee
	Amount         types.AssetAmount  `json:"amount"`
	BlindingFactor types.FixedBuffer  `json:"blinding_factor"`
	From           types.AccountID    `json:"from"`
	Outputs        types.BlindOutputs `json:"outputs"`
}

func (p TransferToBlindOperation) Type() types.OperationType {
	return types.OperationTypeTransferToBlind
}

func (p TransferToBlindOperation) MarshalFeeScheduleParams(params types.M, enc *util.TypeEncoder) error {
	if fee, ok := params["fee"]; ok {
		if err := enc.Encode(types.UInt64(fee.(float64))); err != nil {
			return errors.Annotate(err, "encode Fee")
		}
	}

	if ppk, ok := params["price_per_output"]; ok {
		if err := enc.Encode(types.UInt32(ppk.(float64))); err != nil {
			return errors.Annotate(err, "encode PricePerOutput")
		}
	}

	return nil
}

func (p TransferToBlindOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}
	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}
	if err := enc.Encode(p.Amount); err != nil {
		return errors.Annotate(err, "encode Amount")
	}
	if err := enc.Encode(p.From); err != nil {
		return errors.Annotate(err, "encode From")
	}
	if err := enc.Encode(p.BlindingFactor); err != nil {
		return errors.Annotate(err, "encode BlindingFactor")
	}
	if err := enc.Encode(p.Outputs); err != nil {
		return errors.Annotate(err, "encode Outputs")
	}

	return nil
}
