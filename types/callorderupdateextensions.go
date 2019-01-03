package types

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

//go:generate ffjson $GOFILE

type TargetCollRatio UInt16

func (p TargetCollRatio) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(CallOrderUpdateExtensionsTypeTargetRatio)); err != nil {
		return errors.Annotate(err, "encode type")
	}
	if err := enc.Encode(UInt16(p)); err != nil {
		return errors.Annotate(err, "encode value")
	}

	return nil
}

type CallOrderUpdateExtensions struct {
	TargetCollateralRatio *TargetCollRatio `json:"target_collateral_ratio,omitempty"`
}

func (p CallOrderUpdateExtensions) Length() int {
	fields := 0
	if p.TargetCollateralRatio != nil {
		fields++
	}

	return fields
}

func (p CallOrderUpdateExtensions) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(p.Length())); err != nil {
		return errors.Annotate(err, "encode length")
	}

	if p.TargetCollateralRatio != nil {
		if err := enc.Encode(p.TargetCollateralRatio); err != nil {
			return errors.Annotate(err, "encode TargetCollateralRatio")
		}
	}

	return nil
}
