package types

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type CCDVestingPolicy struct {
	StartClaim     Time   `json:"start_claim"`
	VestingSeconds UInt64 `json:"vesting_seconds"`
}

// TODO: is this yet implemented? test fails!
func (p CCDVestingPolicy) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.StartClaim); err != nil {
		return errors.Annotate(err, "encode StartClaim")
	}

	if err := enc.Encode(p.VestingSeconds); err != nil {
		return errors.Annotate(err, "encode VestingSeconds")
	}

	return nil
}

// test passes
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

type VestingPolicy map[VestingPolicyType]util.TypeMarshaller

func (p *VestingPolicy) UnmarshalJSON(data []byte) error {
	var res []interface{}
	if err := ffjson.Unmarshal(data, &res); err != nil {
		return errors.Annotate(err, "Unmarshal")
	}

	if len(res) != 2 {
		return ErrInvalidInputLength
	}

	(*p) = make(map[VestingPolicyType]util.TypeMarshaller)
	typ := VestingPolicyType(res[0].(float64))

	switch typ {
	case VestingPolicyTypeLinear:
		pol := LinearVestingPolicy{}
		if err := ffjson.Unmarshal(util.ToBytes(res[1]), &pol); err != nil {
			return errors.Annotate(err, "unmarshal LinearVestingPolicy")
		}
		(*p)[0] = pol

	case VestingPolicyTypeCCD:
		pol := CCDVestingPolicy{}
		if err := ffjson.Unmarshal(util.ToBytes(res[1]), &pol); err != nil {
			return errors.Annotate(err, "unmarshal CCDVestingPolicy")
		}
		(*p)[0] = pol
	}

	return nil
}

func (p VestingPolicy) MarshalJSON() ([]byte, error) {
	for k, v := range p {
		return ffjson.Marshal([]interface{}{k, v})
	}

	return nil, nil
}

func (p VestingPolicy) Marshal(enc *util.TypeEncoder) error {
	for k, v := range p {
		if err := enc.EncodeUVarint(uint64(k)); err != nil {
			return errors.Annotate(err, "encode PolicyType")
		}

		if err := enc.Encode(v); err != nil {
			return errors.Annotate(err, "encode Policy")
		}
	}

	return nil
}
