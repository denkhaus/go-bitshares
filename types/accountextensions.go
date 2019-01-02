package types

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type NullExtension struct{}

func (p NullExtension) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(uint8(AccountCreateExtensionsNullExt)); err != nil {
		return errors.Annotate(err, "encode AccountCreateExtensionsNullExt")
	}

	return nil
}

type BuybackOptions struct {
	AssetToBuy       AssetID   `json:"asset_to_buy"`
	AssetToBuyIssuer AccountID `json:"asset_to_buy_issuer"`
	Markets          AssetIDs  `json:"markets"`
}

func (p BuybackOptions) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(uint8(AccountCreateExtensionsBuyback)); err != nil {
		return errors.Annotate(err, "encode AccountCreateExtensionsBuyback")
	}

	if err := enc.Encode(p.AssetToBuy); err != nil {
		return errors.Annotate(err, "encode AssetToBuy")
	}

	if err := enc.Encode(p.AssetToBuyIssuer); err != nil {
		return errors.Annotate(err, "encode AssetToBuyIssuer")
	}

	if err := enc.Encode(p.Markets); err != nil {
		return errors.Annotate(err, "encode Markets")
	}

	return nil
}

type AccountCreateExtensions struct {
	NullExt                *NullExtension          `json:"null_ext,omitempty"`
	OwnerSpecialAuthority  *OwnerSpecialAuthority  `json:"owner_special_authority,omitempty"`
	ActiveSpecialAuthority *ActiveSpecialAuthority `json:"active_special_authority,omitempty"`
	BuybackOptions         *BuybackOptions         `json:"buyback_options,omitempty"`
}

func (p AccountCreateExtensions) Length() int {
	fields := 0
	if p.NullExt != nil {
		fields++
	}
	if p.OwnerSpecialAuthority != nil {
		fields++
	}
	if p.ActiveSpecialAuthority != nil {
		fields++
	}
	if p.BuybackOptions != nil {
		fields++
	}

	return fields
}

func (p AccountCreateExtensions) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(p.Length())); err != nil {
		return errors.Annotate(err, "encode length")
	}

	if p.NullExt != nil {
		if err := enc.Encode(p.NullExt); err != nil {
			return errors.Annotate(err, "encode NullExt")
		}
	}

	if p.OwnerSpecialAuthority != nil {
		if err := enc.Encode(p.OwnerSpecialAuthority); err != nil {
			return errors.Annotate(err, "encode OwnerSpecialAuthority")
		}
	}

	if p.ActiveSpecialAuthority != nil {
		if err := enc.Encode(p.ActiveSpecialAuthority); err != nil {
			return errors.Annotate(err, "encode ActiveSpecialAuthority")
		}
	}

	if p.BuybackOptions != nil {
		if err := enc.Encode(p.BuybackOptions); err != nil {
			return errors.Annotate(err, "encode BuybackOptions")
		}
	}

	return nil
}

type AccountUpdateExtensions struct {
	NullExt                *NullExtension          `json:"null_ext,omitempty"`
	OwnerSpecialAuthority  *OwnerSpecialAuthority  `json:"owner_special_authority,omitempty"`
	ActiveSpecialAuthority *ActiveSpecialAuthority `json:"active_special_authority,omitempty"`
}

func (p AccountUpdateExtensions) Length() int {
	fields := 0
	if p.NullExt != nil {
		fields++
	}
	if p.OwnerSpecialAuthority != nil {
		fields++
	}
	if p.ActiveSpecialAuthority != nil {
		fields++
	}

	return fields
}

func (p AccountUpdateExtensions) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(p.Length())); err != nil {
		return errors.Annotate(err, "encode length")
	}

	if p.NullExt != nil {
		if err := enc.Encode(p.NullExt); err != nil {
			return errors.Annotate(err, "encode NullExt")
		}
	}

	if p.OwnerSpecialAuthority != nil {
		if err := enc.Encode(p.OwnerSpecialAuthority); err != nil {
			return errors.Annotate(err, "encode OwnerSpecialAuthority")
		}
	}

	if p.ActiveSpecialAuthority != nil {
		if err := enc.Encode(p.ActiveSpecialAuthority); err != nil {
			return errors.Annotate(err, "encode ActiveSpecialAuthority")
		}
	}

	return nil
}
