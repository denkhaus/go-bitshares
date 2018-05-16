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

type KeyAuthsMap map[PublicKey]UInt16

func (p *KeyAuthsMap) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		return ErrInvalidInputLength
	}

	var res interface{}
	if err := ffjson.Unmarshal(data, &res); err != nil {
		return errors.Annotate(err, "unmarshal Auths")
	}

	(*p) = make(map[PublicKey]UInt16)
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

		(*p)[PublicKey{key}] = UInt16(weight)
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
		keys = append(keys, k.String())
	}

	sort.Sort(keys, sort.StringComparator)

	for _, k := range keys {
		pub := PublicKey{key: k.(string)}
		if err := enc.Encode(pub); err != nil {
			return errors.Annotate(err, "encode Key")
		}
		if err := enc.Encode(p[pub]); err != nil {
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
		return errors.Annotate(err, "unmarshal Auths")
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
		return nil, errors.Annotate(err, "Marshal")
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

// (string) (len=5) "owner": (map[string]interface {}) (len=4) {
// 	(string) (len=13) "account_auths": ([]interface {}) (len=5 cap=8) {
// 	 ([]interface {}) (len=2 cap=2) {
// 	  (string) (len=10) "1.2.149047",
// 	  (float64) 5
// 	 },
// 	 ([]interface {}) (len=2 cap=2) {
// 	  (string) (len=10) "1.2.386568",
// 	  (float64) 3
// 	 },
// 	 ([]interface {}) (len=2 cap=2) {
// 	  (string) (len=10) "1.2.386686",
// 	  (float64) 4
// 	 },
// 	 ([]interface {}) (len=2 cap=2) {
// 	  (string) (len=10) "1.2.395052",
// 	  (float64) 4
// 	 },
// 	 ([]interface {}) (len=2 cap=2) {
// 	  (string) (len=10) "1.2.442608",
// 	  (float64) 3
// 	 }
// 	},
// 	(string) (len=13) "address_auths": ([]interface {}) {
// 	},
// 	(string) (len=9) "key_auths": ([]interface {}) {
// 	},
// 	(string) (len=16) "weight_threshold": (float64) 10
//    }

type AccountAuthsMap map[GrapheneID]UInt16

func (p *AccountAuthsMap) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		return nil
	}

	var res interface{}
	if err := ffjson.Unmarshal(data, &res); err != nil {
		return errors.Annotate(err, "unmarshal Auths")
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
		return nil, errors.Annotate(err, "Marshal")
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

type SpecialAuthsMap map[string]interface{}

func (p *SpecialAuthsMap) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		return nil
	}

	var res interface{}
	if err := ffjson.Unmarshal(data, &res); err != nil {
		return errors.Annotate(err, "unmarshal Auths")
	}

	// (*p) = make(map[string]int64)
	// auths := res.([]interface{})

	// for _, a := range auths {
	// 	tk := a.([]interface{})
	// 	(*p)[tk[0].(string)] = int64(tk[1].(float64))
	// }

	return nil
}
