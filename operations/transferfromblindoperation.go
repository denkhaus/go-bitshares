package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeTransferFromBlind] = func() types.Operation {
		op := &TransferFromBlindOperation{}
		return op
	}
}

type TransferFromBlindOperation struct {
	types.OperationFee
	Amount      types.AssetAmount `json:"amount"`
	To          types.GrapheneID  `json:"to"`
	BlindFactor types.FixedBuffer `json:"blinding_factor"`
	BlindInputs types.BlindInputs `json:"inputs"`
}

func (p TransferFromBlindOperation) Type() types.OperationType {
	return types.OperationTypeTransferFromBlind
}

func (p TransferFromBlindOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}
	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}
	if err := enc.Encode(p.Amount); err != nil {
		return errors.Annotate(err, "encode Amount")
	}
	if err := enc.Encode(p.To); err != nil {
		return errors.Annotate(err, "encode To")
	}
	if err := enc.Encode(p.BlindFactor); err != nil {
		return errors.Annotate(err, "encode BlindFactor")
	}
	if err := enc.Encode(p.BlindInputs); err != nil {
		return errors.Annotate(err, "encode BlindInputs")
	}

	return nil
}
