package types

//go:generate ffjson   $GOFILE

import (
	"github.com/denkhaus/bitshares/util"
	sort "github.com/emirpasic/gods/utils"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

//TODO implement AddressAuths
type Authority struct {
	WeightThreshold UInt32          `json:"weight_threshold"`
	AccountAuths    AccountAuthsMap `json:"account_auths"`
	KeyAuths        KeyAuthsMap     `json:"key_auths"`
	AddressAuths    AuthsMap        `json:"address_auths"`
	Extensions      Extensions      `json:"extensions"`
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

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extensions")
	}

	return nil
}

type KeyAuthsMap map[string]UInt16

func (p *KeyAuthsMap) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		return ErrInvalidInputLength
	}

	var res interface{}
	if err := ffjson.Unmarshal(data, &res); err != nil {
		return errors.Annotate(err, "unmarshal Auths")
	}

	(*p) = make(map[string]UInt16)
	auths, ok := res.([]interface{})
	if !ok {
		return ErrInvalidInputType
	}

	for _, a := range auths {
		tk, ok := a.([]interface{})
		if !ok {
			return ErrInvalidInputType
		}

		key, ok := tk[0].(string)
		if !ok {
			return ErrInvalidInputType
		}

		weight, ok := tk[1].(float64)
		if !ok {
			return ErrInvalidInputType
		}

		pub, err := NewPublicKey(key)
		if err != nil {
			return errors.Annotate(err, "NewPublicKey")
		}

		(*p)[pub.String()] = UInt16(weight)
	}

	return nil
}

func (p KeyAuthsMap) MarshalJSON() ([]byte, error) {
	ret := []interface{}{}

	for k, v := range p {
		ret = append(ret, []interface{}{k, v})
	}

	buf, err := ffjson.Marshal(ret)
	if err != nil {
		return nil, errors.Annotate(err, "Marshal")
	}

	return buf, nil
}

func (p KeyAuthsMap) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	// copy keys
	keys := []interface{}{}
	for k := range p {
		keys = append(keys, k)
	}

	sort.Sort(keys, sort.StringComparator)

	for _, k := range keys {
		key := k.(string)
		pub, err := NewPublicKey(key)
		if err != nil {
			return errors.Annotate(err, "NewPublicKey")
		}

		if err := pub.Marshal(enc); err != nil {
			return errors.Annotate(err, "encode Key")
		}

		if err := enc.Encode(p[key]); err != nil {
			return errors.Annotate(err, "encode ValueExtension")
		}
	}

	return nil
}

type AuthsMap map[string]UInt16

func (p *AuthsMap) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		return nil
	}

	var res interface{}
	if err := ffjson.Unmarshal(data, &res); err != nil {
		return errors.Annotate(err, "unmarshal AuthsMap")
	}

	(*p) = make(map[string]UInt16)
	auths := res.([]interface{})

	for _, a := range auths {
		tk := a.([]interface{})
		(*p)[tk[0].(string)] = UInt16(tk[1].(float64))
	}

	return nil
}

func (p AuthsMap) MarshalJSON() ([]byte, error) {
	ret := []interface{}{}

	for k, v := range p {
		ret = append(ret, []interface{}{k, v})
	}

	buf, err := ffjson.Marshal(ret)
	if err != nil {
		return nil, errors.Annotate(err, "marshal AuthsMap")
	}

	return buf, nil
}

func (p AuthsMap) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	for k, v := range p {
		if err := enc.Encode(k); err != nil {
			return errors.Annotate(err, "encode Key")
		}
		if err := enc.Encode(v); err != nil {
			return errors.Annotate(err, "encode ValueExtension")
		}
	}

	return nil
}

type AccountAuthsMap map[GrapheneID]UInt16

func (p *AccountAuthsMap) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		return nil
	}

	var res interface{}
	if err := ffjson.Unmarshal(data, &res); err != nil {
		return errors.Annotate(err, "unmarshal AccountAuthsMap")
	}

	(*p) = make(map[GrapheneID]UInt16)
	auths := res.([]interface{})

	for _, a := range auths {
		tk := a.([]interface{})
		(*p)[*NewGrapheneID(tk[0].(string))] = UInt16(tk[1].(float64))
	}

	return nil
}

func (p AccountAuthsMap) MarshalJSON() ([]byte, error) {
	ret := []interface{}{}

	for k, v := range p {
		ret = append(ret, []interface{}{k, v})
	}

	buf, err := ffjson.Marshal(ret)
	if err != nil {
		return nil, errors.Annotate(err, "marshal AccountAuthsMap")
	}

	return buf, nil
}

func (p AccountAuthsMap) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	for k, v := range p {
		if err := enc.Encode(k); err != nil {
			return errors.Annotate(err, "encode Key")
		}
		if err := enc.Encode(v); err != nil {
			return errors.Annotate(err, "encode Value")
		}
	}

	return nil
}

type NoSpecialAuthority map[string]interface{}

type TopHoldersSpecialAuthority struct {
	Asset         GrapheneID `json:"asset"`
	NumTopHolders UInt8      `json:"num_top_holders"`
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
	typ  SpecialAuthorityType
	data interface{}
}

func (p *SpecialAuth) UnmarshalJSON(data []byte) error {
	d := []interface{}{}
	if err := ffjson.Unmarshal(data, &d); err != nil {
		return errors.Annotate(err, "unmarshal SpecialAuthsMap")
	}

	if len(d) != 2 {
		return ErrInvalidInputLength
	}

	t, ok := d[0].(float64)
	if !ok {
		return ErrInvalidInputType
	}

	p.typ = SpecialAuthorityType(t)

	switch p.typ {
	case SpecialAuthorityTypeNoSpecial:
		p.data = &NoSpecialAuthority{}
	case SpecialAuthorityTypeTopHolders:
		p.data = &TopHoldersSpecialAuthority{}
	}

	if err := ffjson.Unmarshal(util.ToBytes(d[1]), p.data); err != nil {
		return errors.Annotate(err, "unmarshal SpecialAuthority")
	}

	return nil
}

func (p SpecialAuth) MarshalJSON() ([]byte, error) {
	return ffjson.Marshal([]interface{}{
		p.typ,
		p.data,
	})
}

func (p SpecialAuth) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(uint8(p.typ)); err != nil {
		return errors.Annotate(err, "encode Type")
	}

	switch p.typ {
	case SpecialAuthorityTypeNoSpecial:
		if err := enc.Encode(p.data.(*NoSpecialAuthority)); err != nil {
			return errors.Annotate(err, "encode Data")
		}
	case SpecialAuthorityTypeTopHolders:
		if err := enc.Encode(p.data.(*TopHoldersSpecialAuthority)); err != nil {
			return errors.Annotate(err, "encode Data")
		}
	}

	return nil
}
