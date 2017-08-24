package objects

import (
	json "encoding/json"

	"github.com/juju/errors"
)

type Authority struct {
	WeightThreshold UInt64          `json:"weight_threshold"`
	AccountAuths    MapAccountAuths `json:"account_auths"`
	KeyAuths        MapKeyAuths     `json:"key_auths"`
	Extensions      Extensions      `json:"extensions"`
}

type MapAccountAuths map[GrapheneID]int64

func (p *MapAccountAuths) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		return nil
	}

	var res interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return errors.Annotate(err, "unmarshal MapAccountAuths")
	}

	(*p) = make(map[GrapheneID]int64)
	accAuths := res.([]interface{})

	for _, aa := range accAuths {
		tk := aa.([]interface{})
		(*p)[*NewGrapheneID(ObjectID(tk[0].(string)))] = int64(tk[1].(float64))
	}

	return nil
}

type MapKeyAuths map[PublicKey]int64

func (p *MapKeyAuths) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		return nil
	}

	var res interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return errors.Annotate(err, "unmarshal MapKeyAuths")
	}

	(*p) = make(map[PublicKey]int64)
	keyAuths := res.([]interface{})

	for _, ka := range keyAuths {
		tk := ka.([]interface{})
		(*p)[*NewPublicKey(tk[0].(string))] = int64(tk[1].(float64))
	}

	return nil
}
