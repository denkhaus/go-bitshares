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

	return nil
}
