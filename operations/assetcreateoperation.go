package operations

//go:generate ffjson $GOFILE

import (
	"strconv"

	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeAssetCreate] = func() types.Operation {
		op := &AssetCreateOperation{}
		return op
	}
}

type AssetCreateOperation struct {
	types.OperationFee
	BitassetOptions    *types.BitassetOptions `json:"bitasset_opts"`
	CommonOptions      types.AssetOptions     `json:"common_options"`
	Extensions         types.Extensions       `json:"extensions"`
	IsPredictionMarket bool                   `json:"is_prediction_market"`
	Issuer             types.GrapheneID       `json:"issuer"`
	Precision          types.UInt8            `json:"precision"`
	Symbol             string                 `json:"symbol"`
}

func (p AssetCreateOperation) Type() types.OperationType {
	return types.OperationTypeAssetCreate
}

func (p AssetCreateOperation) MarshalFeeScheduleParams(params types.M, enc *util.TypeEncoder) error {
	if s3, ok := params["symbol3"]; ok {
		s3Val, err := ToUInt64(s3)
		if err != nil {
			return errors.Annotate(err, "ToUint64 [symbol3]")
		}
		if err := enc.Encode(s3Val); err != nil {
			return errors.Annotate(err, "encode Symbol3")
		}
	}
	if s4, ok := params["symbol4"]; ok {
		s4Val, err := ToUInt64(s4)
		if err != nil {
			return errors.Annotate(err, "ToUint64 [symbol4]")
		}
		if err := enc.Encode(s4Val); err != nil {
			return errors.Annotate(err, "encode Symbol4")
		}
	}
	if ls, ok := params["long_symbol"]; ok {
		lsVal, err := ToUInt64(ls)
		if err != nil {
			return errors.Annotate(err, "ToUint64 [LongSymbol]")
		}
		if err := enc.Encode(lsVal); err != nil {
			return errors.Annotate(err, "encode LongSymbol")
		}
	}
	if ppk, ok := params["price_per_kbyte"]; ok {
		if err := enc.Encode(types.UInt32(ppk.(float64))); err != nil {
			return errors.Annotate(err, "encode PricePerKByte")
		}
	}

	return nil
}

func (p AssetCreateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Issuer); err != nil {
		return errors.Annotate(err, "encode issuer")
	}

	if err := enc.Encode(p.Symbol); err != nil {
		return errors.Annotate(err, "encode Symbol")
	}

	if err := enc.Encode(p.Precision); err != nil {
		return errors.Annotate(err, "encode Precision")
	}

	if err := enc.Encode(p.CommonOptions); err != nil {
		return errors.Annotate(err, "encode CommonOptions")
	}

	if err := enc.Encode(p.BitassetOptions != nil); err != nil {
		return errors.Annotate(err, "encode have BitassetOptions")
	}

	if err := enc.Encode(p.BitassetOptions); err != nil {
		return errors.Annotate(err, "encode BitassetOptions")
	}

	if err := enc.Encode(p.IsPredictionMarket); err != nil {
		return errors.Annotate(err, "encode IsPredictionMarket")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

func ToUInt64(in interface{}) (types.UInt64, error) {
	inString, ok := in.(string)
	if ok {
		inVal, err := strconv.ParseUint(inString, 10, 64)
		if err != nil {
			return types.UInt64(0), errors.Annotate(err, "ParseUint")
		}
		return types.UInt64(inVal), nil
	}

	return types.UInt64(in.(float64)), nil
}
