package types

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

//TODO: CoinSecondsEarned is UInt128! Since golang has no
//128 bit uint, check marshal and implement this later.
type CCDVestingPolicy struct {
	StartClaim                  Time   `json:"start_claim"`
	CoinSecondsEarnedLastUpdate Time   `json:"coin_seconds_earned_last_update"`
	VestingSeconds              UInt32 `json:"vesting_seconds"`
	CoinSecondsEarned           UInt64 `json:"coin_seconds_earned"` //UInt128!!
}

// TODO: check order!
func (p CCDVestingPolicy) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.StartClaim); err != nil {
		return errors.Annotate(err, "encode StartClaim")
	}

	if err := enc.Encode(p.VestingSeconds); err != nil {
		return errors.Annotate(err, "encode VestingSeconds")
	}

	// if err := enc.Encode(p.CoinSecondsEarnedLastUpdate); err != nil {
	// 	return errors.Annotate(err, "encode CoinSecondsEarnedLastUpdate")
	// }

	// if err := enc.Encode(p.CoinSecondsEarned); err != nil {
	// 	return errors.Annotate(err, "encode CoinSecondsEarned")
	// }

	return nil
}

type LinearVestingPolicy struct {
	BeginTimestamp         Time   `json:"begin_timestamp"`
	VestingCliffSeconds    UInt32 `json:"vesting_cliff_seconds"`
	VestingDurationSeconds UInt32 `json:"vesting_duration_seconds"`
}

func (p LinearVestingPolicy) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.BeginTimestamp); err != nil {
		return errors.Annotate(err, "encode BeginTimestamp")
	}

	if err := enc.Encode(p.VestingCliffSeconds); err != nil {
		return errors.Annotate(err, "encode VestingCliffSeconds")
	}

	if err := enc.Encode(p.VestingDurationSeconds); err != nil {
		return errors.Annotate(err, "encode VestingDurationSeconds")
	}

	return nil
}

type VestingPolicy struct {
	typ  VestingPolicyType
	data util.TypeMarshaler
}

func (p *VestingPolicy) UnmarshalJSON(data []byte) error {
	var res []interface{}
	if err := ffjson.Unmarshal(data, &res); err != nil {
		return errors.Annotate(err, "Unmarshal")
	}

	if len(res) != 2 {
		return ErrInvalidInputLength
	}

	p.typ = VestingPolicyType(res[0].(float64))

	switch p.typ {
	case VestingPolicyTypeLinear:
		pol := LinearVestingPolicy{}
		if err := ffjson.Unmarshal(util.ToBytes(res[1]), &pol); err != nil {
			return errors.Annotate(err, "unmarshal LinearVestingPolicy")
		}
		p.data = pol

	case VestingPolicyTypeCCD:
		pol := CCDVestingPolicy{}
		if err := ffjson.Unmarshal(util.ToBytes(res[1]), &pol); err != nil {
			return errors.Annotate(err, "unmarshal CCDVestingPolicy")
		}
		p.data = pol
	}

	return nil
}

func (p VestingPolicy) MarshalJSON() ([]byte, error) {
	return ffjson.Marshal([]interface{}{
		p.typ,
		p.data,
	})
}

func (p VestingPolicy) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(p.typ)); err != nil {
		return errors.Annotate(err, "encode PolicyType")
	}

	if err := enc.Encode(p.data); err != nil {
		return errors.Annotate(err, "encode Policy")
	}

	return nil
}
