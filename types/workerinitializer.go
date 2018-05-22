package types

import (
	"encoding/json"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type RefundWorkerInitializer struct {
}

func (p RefundWorkerInitializer) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(uint8(WorkerInitializerTypeRefund)); err != nil {
		return errors.Annotate(err, "encode Type")
	}

	return nil
}

type VestingBalanceWorkerInitializer struct {
	PayVestingPeriodDays UInt16 `json:"pay_vesting_period_days"`
}

func (p VestingBalanceWorkerInitializer) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(uint8(WorkerInitializerTypeVestingBalance)); err != nil {
		return errors.Annotate(err, "encode Type")
	}

	if err := enc.Encode(p.PayVestingPeriodDays); err != nil {
		return errors.Annotate(err, "encode PayVestingPeriodDays")
	}

	return nil
}

type BurnWorkerInitializer struct {
}

func (p BurnWorkerInitializer) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(uint8(WorkerInitializerTypeBurn)); err != nil {
		return errors.Annotate(err, "encode Type")
	}

	return nil
}

type WorkerInitializer struct {
	Type        WorkerInitializerType
	Initializer interface{}
}

func (p WorkerInitializer) Marshal(enc *util.TypeEncoder) error {
	switch p.Type {
	case WorkerInitializerTypeRefund:
		if err := enc.Encode(p.Initializer.(*RefundWorkerInitializer)); err != nil {
			return errors.Annotate(err, "encode Initializer")
		}
	case WorkerInitializerTypeVestingBalance:
		if err := enc.Encode(p.Initializer.(*VestingBalanceWorkerInitializer)); err != nil {
			return errors.Annotate(err, "encode Initializer")
		}
	case WorkerInitializerTypeBurn:
		if err := enc.Encode(p.Initializer.(*BurnWorkerInitializer)); err != nil {
			return errors.Annotate(err, "encode Initializer")
		}
	}

	return nil
}

func (p WorkerInitializer) MarshalJSON() ([]byte, error) {
	return ffjson.Marshal([]interface{}{
		p.Type,
		p.Initializer,
	})
}

func (p *WorkerInitializer) UnmarshalJSON(data []byte) error {
	raw := make([]json.RawMessage, 2)
	if err := ffjson.Unmarshal(data, &raw); err != nil {
		return errors.Annotate(err, "unmarshal RawMessage")
	}

	if len(raw) != 2 {
		return ErrInvalidInputLength
	}

	if err := ffjson.Unmarshal(raw[0], &p.Type); err != nil {
		return errors.Annotate(err, "unmarshal InitializerType")
	}

	switch p.Type {
	case WorkerInitializerTypeRefund:
		p.Initializer = &RefundWorkerInitializer{}
	case WorkerInitializerTypeVestingBalance:
		p.Initializer = &VestingBalanceWorkerInitializer{}
	case WorkerInitializerTypeBurn:
		p.Initializer = &BurnWorkerInitializer{}
	}

	if err := ffjson.Unmarshal(raw[1], p.Initializer); err != nil {
		return errors.Annotate(err, "unmarshal Initializer")
	}

	return nil
}
