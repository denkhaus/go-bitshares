package types

import (
	"encoding/json"

	"github.com/denkhaus/bitshares/util"
	"github.com/denkhaus/logging"
	sort "github.com/emirpasic/gods/utils"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type FeeScheduleParameter struct {
	OperationType OperationType
	Params        M
}

func (p FeeScheduleParameter) MarshalJSON() ([]byte, error) {
	return ffjson.Marshal([]interface{}{
		p.OperationType,
		p.Params,
	})
}

func (p *FeeScheduleParameter) UnmarshalJSON(data []byte) error {
	raw := make([]json.RawMessage, 2)
	if err := ffjson.Unmarshal(data, &raw); err != nil {
		return errors.Annotate(err, "unmarshal RawData")
	}

	if err := ffjson.Unmarshal(raw[0], &p.OperationType); err != nil {
		return errors.Annotate(err, "unmarshal OperationType")
	}

	if err := ffjson.Unmarshal(raw[1], &p.Params); err != nil {
		return errors.Annotate(err, "unmarshal Params")
	}

	return nil
}

func (p FeeScheduleParameter) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(uint8(p.OperationType)); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if getOp, ok := OperationMap[p.OperationType]; ok {
		if err := getOp().MarshalFeeScheduleParams(p.Params, enc); err != nil {
			return errors.Annotate(err, "MarshalFeeScheduleParams")
		}
	} else {
		logging.Warnf("FeeSchedule marshaling is not yet implemented for %s",
			p.OperationType.OperationName())
	}

	return nil
}

type FeeScheduleParameters []FeeScheduleParameter

func (p FeeScheduleParameters) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	params := make([]interface{}, 0, len(p))
	for _, pa := range p {
		params = append(params, pa)
	}

	sort.Sort(params, func(a, b interface{}) int {
		aParam := a.(FeeScheduleParameter)
		bParam := b.(FeeScheduleParameter)

		return sort.Int8Comparator(
			int8(aParam.OperationType),
			int8(bParam.OperationType),
		)
	})

	for _, p := range params {
		if err := enc.Encode(p.(FeeScheduleParameter)); err != nil {
			return errors.Annotate(err, "encode Parameter")
		}
	}

	return nil
}

type FeeSchedule struct {
	Scale      UInt32                `json:"scale"`
	Parameters FeeScheduleParameters `json:"parameters"`
}

func (p FeeSchedule) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.Parameters); err != nil {
		return errors.Annotate(err, "encode Parameters")
	}

	if err := enc.Encode(p.Scale); err != nil {
		return errors.Annotate(err, "encode Scale")
	}

	return nil
}
