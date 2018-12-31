package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/denkhaus/logging"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeCustom] = func() types.Operation {
		op := &CustomOperation{}
		return op
	}
}

type CustomOperation struct {
	types.OperationFee

	// asset                     fee;
	// account_id_type           payer;
	// flat_set<account_id_type> required_auths;
	// uint16_t                  id = 0;
	// vector<char> data;
}

func (p CustomOperation) Type() types.OperationType {
	return types.OperationTypeCustom
}

func (p CustomOperation) MarshalFeeScheduleParams(params types.M, enc *util.TypeEncoder) error {
	if fee, ok := params["fee"]; ok {
		if err := enc.Encode(types.UInt64(fee.(float64))); err != nil {
			return errors.Annotate(err, "encode Fee")
		}
	}

	if ppk, ok := params["price_per_kbyte"]; ok {
		if err := enc.Encode(types.UInt32(ppk.(float64))); err != nil {
			return errors.Annotate(err, "encode PricePerKByte")
		}
	}

	return nil
}

func (p CustomOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}
	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}
	//(fee)(payer)(required_auths)(id)(data)
	logging.Warnf("%s is not implemented", p.Type().OperationName())
	return nil
}
