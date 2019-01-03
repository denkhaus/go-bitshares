package types

//go:generate ffjson $GOFILE

import (
	"encoding/json"

	"github.com/denkhaus/bitshares/util"
	sort "github.com/emirpasic/gods/utils"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type Authority struct {
	WeightThreshold UInt32          `json:"weight_threshold"`
	AccountAuths    AccountAuthsMap `json:"account_auths"`
	KeyAuths        KeyAuthsMap     `json:"key_auths"`
	AddressAuths    AddressAuthsMap `json:"address_auths"`
}

func (p Authority) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.WeightThreshold); err != nil {
		return errors.Annotate(err, "encode WeightThreshold")
	}
	if err := enc.Encode(p.AccountAuths); err != nil {
		return errors.Annotate(err, "encode AccountAuths")
	}
	if err := enc.Encode(p.KeyAuths); err != nil {
		return errors.Annotate(err, "encode KeyAuths")
	}
	if err := enc.Encode(p.AddressAuths); err != nil {
		return errors.Annotate(err, "encode AddressAuths")
	}

	return nil
}

type KeyAuthsMap map[*PublicKey]UInt16

func (p *KeyAuthsMap) UnmarshalJSON(data []byte) error {
	var auths [][]interface{}
	if err := ffjson.Unmarshal(data, &auths); err != nil {
		return errors.Annotate(err, "unmarshal Auths")
	}

	(*p) = make(map[*PublicKey]UInt16)
	for _, tk := range auths {
		key, ok := tk[0].(string)
		if !ok {
			return ErrInvalidInputType
		}

		weight, ok := tk[1].(float64)
		if !ok {
			return ErrInvalidInputType
		}

		pub, err := NewPublicKeyFromString(key)
		if err != nil {
			return errors.Annotate(err, "NewPublicKeyFromString")
		}

		(*p)[pub] = UInt16(weight)
	}

	return nil
}

func (p KeyAuthsMap) MarshalJSON() ([]byte, error) {
	ret := make([]interface{}, 0, len(p))
	for k, v := range p {
		ret = append(ret, []interface{}{k.String(), v})
	}
	return ffjson.Marshal(ret)
}

func (p KeyAuthsMap) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	//sort keys
	keys := make([]interface{}, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}

	var err error
	sort.Sort(keys, func(a, b interface{}) (s int) {
		s, err = publicKeyComparator(a.(*PublicKey), b.(*PublicKey))
		return
	})

	if err != nil {
		return errors.Annotate(err, "Sort")
	}

	for _, k := range keys {
		pub := k.(*PublicKey)
		if err := pub.Marshal(enc); err != nil {
			return errors.Annotate(err, "encode PubKey")
		}

		if err := enc.Encode(p[pub]); err != nil {
			return errors.Annotate(err, "encode Weight")
		}
	}

	return nil
}

type AddressAuthsMap map[*Address]UInt16

func (p *AddressAuthsMap) UnmarshalJSON(data []byte) error {
	var auths [][]interface{}
	if err := ffjson.Unmarshal(data, &auths); err != nil {
		return errors.Annotate(err, "unmarshal AddressAuthsMap")
	}

	(*p) = make(map[*Address]UInt16)
	for _, tk := range auths {
		add, ok := tk[0].(string)
		if !ok {
			return ErrInvalidInputType
		}
		addr, err := NewAddressFromString(add)
		if err != nil {
			return errors.Annotate(err, "NewAddressFromString")
		}

		weight, ok := tk[1].(float64)
		if !ok {
			return ErrInvalidInputType
		}

		(*p)[addr] = UInt16(weight)
	}

	return nil
}

func (p AddressAuthsMap) MarshalJSON() ([]byte, error) {
	ret := []interface{}{}
	for k, v := range p {
		ret = append(ret, []interface{}{k, v})
	}
	return ffjson.Marshal(ret)
}

func (p AddressAuthsMap) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	//sort keys
	keys := make([]interface{}, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}

	sort.Sort(keys, func(a, b interface{}) (s int) {
		return sort.StringComparator(
			a.(*Address).String(),
			b.(*Address).String(),
		)
	})

	for _, v := range keys {
		add := v.(*Address)
		if err := enc.Encode(add); err != nil {
			return errors.Annotate(err, "encode Address")
		}
		if err := enc.Encode(p[add]); err != nil {
			return errors.Annotate(err, "encode Weight")
		}
	}

	return nil
}

type AccountAuthsMap map[GrapheneObject]UInt16

func (p *AccountAuthsMap) UnmarshalJSON(data []byte) error {
	var auths [][]interface{}
	if err := ffjson.Unmarshal(data, &auths); err != nil {
		return errors.Annotate(err, "Unmarshal")
	}

	(*p) = make(map[GrapheneObject]UInt16)
	for _, tk := range auths {
		acc, ok := tk[0].(string)
		if !ok {
			return ErrInvalidInputType
		}

		weight, ok := tk[1].(float64)
		if !ok {
			return ErrInvalidInputType
		}

		(*p)[NewAccountID(acc)] = UInt16(weight)
	}

	return nil
}

func (p AccountAuthsMap) MarshalJSON() ([]byte, error) {
	ret := []interface{}{}
	for k, v := range p {
		ret = append(ret, []interface{}{k, v})
	}

	return ffjson.Marshal(ret)
}

func (p AccountAuthsMap) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	//sort keys
	keys := make([]interface{}, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}

	sort.Sort(keys, func(a, b interface{}) (s int) {
		return sort.UInt64Comparator(
			uint64(a.(GrapheneObject).Instance()),
			uint64(b.(GrapheneObject).Instance()),
		)
	})

	for _, k := range keys {
		ob := k.(GrapheneObject)
		if err := enc.Encode(ob); err != nil {
			return errors.Annotate(err, "encode Account")
		}

		if err := enc.Encode(p[ob]); err != nil {
			return errors.Annotate(err, "encode Weight")
		}
	}

	return nil
}

type NoSpecialAuthority struct{}

type TopHoldersSpecialAuthority struct {
	Asset         AssetID `json:"asset"`
	NumTopHolders UInt8   `json:"num_top_holders"`
}

func (p TopHoldersSpecialAuthority) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.Asset); err != nil {
		return errors.Annotate(err, "encode Asset")
	}

	if err := enc.Encode(p.NumTopHolders); err != nil {
		return errors.Annotate(err, "encode NumTopHolders")
	}

	return nil
}

type OwnerSpecialAuthority struct {
	SpecialAuth
}

func (p OwnerSpecialAuthority) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(uint8(AccountCreateExtensionsOwnerSpecial)); err != nil {
		return errors.Annotate(err, "encode AccountCreateExtensionsOwnerSpecial")
	}

	if err := enc.Encode(p.SpecialAuth); err != nil {
		return errors.Annotate(err, "encode SpecialAuth")
	}

	return nil
}

type ActiveSpecialAuthority struct {
	SpecialAuth
}

func (p ActiveSpecialAuthority) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(uint8(AccountCreateExtensionsActiveSpecial)); err != nil {
		return errors.Annotate(err, "encode AccountCreateExtensionsActiveSpecial")
	}

	if err := enc.Encode(p.SpecialAuth); err != nil {
		return errors.Annotate(err, "encode SpecialAuths")
	}

	return nil
}

type SpecialAuth struct {
	Type SpecialAuthorityType
	Auth interface{}
}

func (p *SpecialAuth) UnmarshalJSON(data []byte) error {
	raw := make([]json.RawMessage, 2)
	if err := ffjson.Unmarshal(data, &raw); err != nil {
		return errors.Annotate(err, "unmarshal RawData")
	}

	if err := ffjson.Unmarshal(raw[0], &p.Type); err != nil {
		return errors.Annotate(err, "unmarshal AuthorityType")
	}

	switch p.Type {
	case SpecialAuthorityTypeNoSpecial:
		p.Auth = &NoSpecialAuthority{}
	case SpecialAuthorityTypeTopHolders:
		p.Auth = &TopHoldersSpecialAuthority{}
	}

	if err := ffjson.Unmarshal(raw[1], p.Auth); err != nil {
		return errors.Annotate(err, "unmarshal SpecialAuthority")
	}

	return nil
}

func (p SpecialAuth) MarshalJSON() ([]byte, error) {
	return ffjson.Marshal([]interface{}{
		p.Type,
		p.Auth,
	})
}

func (p SpecialAuth) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(uint8(p.Type)); err != nil {
		return errors.Annotate(err, "encode Type")
	}

	switch p.Type {
	case SpecialAuthorityTypeNoSpecial:
		if err := enc.Encode(p.Auth.(*NoSpecialAuthority)); err != nil {
			return errors.Annotate(err, "encode Data")
		}
	case SpecialAuthorityTypeTopHolders:
		if err := enc.Encode(p.Auth.(*TopHoldersSpecialAuthority)); err != nil {
			return errors.Annotate(err, "encode Data")
		}
	}

	return nil
}
