package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/denkhaus/logging"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeAccountTransfer] = func() types.Operation {
		op := &AccountTransferOperation{}
		return op
	}
}

type AccountTransferOperation struct {
	types.OperationFee
}

func (p AccountTransferOperation) Type() types.OperationType {
	return types.OperationTypeAccountTransfer
}

func (p AccountTransferOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	logging.Warnf("%s is not implemented", p.Type().OperationName())
	return nil
}
