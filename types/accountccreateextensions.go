package types

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type NullExtension interface{}

type BuybackOptions struct {
	AssetToBuy       GrapheneID  `json:"asset_to_buy"`
	AssetToBuyIssuer GrapheneID  `json:"asset_to_buy_issuer"`
	Markets          GrapheneIDs `json:"markets"`
}

func (p BuybackOptions) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.AssetToBuy); err != nil {
		return errors.Annotate(err, "encode AssetToBuy")
	}

	if err := enc.Encode(p.AssetToBuyIssuer); err != nil {
		return errors.Annotate(err, "encode AssetToBuyIssuer")
	}

	// if err := enc.Encode(p.Markets); err != nil {
	// 	return errors.Annotate(err, "encode Markets")
	// }

	return nil
}

type AccountCreateExtensions struct {
	BuybackOptions BuybackOptions
	// NullExtension         NullExtension
	// OwnerSpecialAuthority SpecialAuthsMap
}

func (p *AccountCreateExtensions) UnmarshalJSON(data []byte) error {

	res := make(map[string]interface{})

	if err := ffjson.Unmarshal(data, &res); err != nil {
		return errors.Annotate(err, "unmarshal")
	}

	if len(res) > 0 {
		util.DumpJSON("ext", res)
	}

	return nil
}

func (p AccountCreateExtensions) MarshalJSON() ([]byte, error) {
	ret := []interface{}{}

	if p.BuybackOptions.AssetToBuy.Valid() {
		ret = append(ret, []interface{}{
			p.BuybackOptions, 3,
		})
	}

	buf, err := ffjson.Marshal(ret)
	if err != nil {
		return nil, errors.Annotate(err, "marshal AccountCreateExtensions")
	}

	return buf, nil
}

//TODO: define this
func (p AccountCreateExtensions) Marshal(enc *util.TypeEncoder) error {
	// if err := enc.Encode(p.BuybackOptions); err != nil {
	// 	return errors.Annotate(err, "encode BuybackOptions")
	// }

	return nil
}

// `json:"buyback_options,omitempty"`
// `json:"null_ext,omitempty"`
// `json:"owner_special_authority,omitempty"`
