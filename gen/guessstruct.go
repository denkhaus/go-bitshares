package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/denkhaus/bitshares/gen/data"
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/denkhaus/gojson"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

var (
	// do not change order here
	knownTypes = []types.Unmarshalable{
		//&types.AccountOptions{},
		// &types.Asset{},
		&types.AssetAmount{},
		&types.AssetFeed{},
		// &types.AssetOptions{},
		&types.ObjectID{},
		//&types.BitAssetDataOptions{},
		&types.Authority{},
		//&types.Memo{},
		&types.Price{},
		&types.PriceFeed{},
		//&types.Votes{},
		&types.Time{},
		//&types.PublicKey{},
		//&types.Account{},
	}
)

func generateOpData(d GenData) error {
	samples, err := data.GetSamplesByType(d.Type)
	if err != nil {
		return errors.Annotate(err, "GetSampleByType")
	}

	for _, s := range samples {
		sample, err := strconv.Unquote(s)
		if err != nil {
			return errors.Annotate(err, "Unquote")
		}

		//fmt.Printf("generate struct by sample %+v\n", sample)

		buf, err := gojson.GenerateWithTypeGuessing(
			strings.NewReader(sample),
			gojson.ParseJson, d.Type.OperationName(),
			"operations", []string{"json"}, true, true,
			guessStructType,
		)

		if err != nil {
			return errors.Annotate(err, "GenerateWithTypeGuessing")
		}

		fmt.Println("generated struct ", string(buf))
	}

	return nil
}

func guessStructType(value interface{}, suggestedType string) (string, error) {
	//util.Dump("valueToGuess", value)
	//	util.DumpJSON("suggestedType", suggestedType)

	for _, t := range knownTypes {
		v, err := ffjson.Marshal(value)
		if err != nil {
			return "", errors.Annotate(err, "Marshal")
		}

		//make local copy of known type
		typ := t

		//	util.Dump("data", v)
		if err := typ.UnmarshalJSON(v); err == nil {
			// util.Dump("compare-1", typ)
			// util.Dump("compare-2", value)

			switch o := typ.(type) {
			case *types.ObjectID:
				if value.(string) == o.String() {
					return "types.ObjectID", nil
				}
			case *types.Time:
				if value.(string) == o.String() {
					return "types.Time", nil
				}
			// case *types.PublicKey:
			// 	if value.(string) == o.String() {
			// 		return "types.PublicKey", nil
			// 	}
			case *types.AccountOptions:
				if bytes.Equal(v, util.ToBytes(typ)) {
					return "types.AccountOptions", nil
				}
			case *types.AssetAmount:
				if bytes.Equal(v, util.ToBytes(typ)) {
					return "types.AssetAmount", nil
				}
			case *types.Authority:
				if bytes.Equal(v, util.ToBytes(typ)) {
					return "types.Authority", nil
				}
			case *types.Memo:
				if bytes.Equal(v, util.ToBytes(typ)) {
					return "types.Memo", nil
				}
			// case *types.Votes:
			// 	if bytes.Equal(v, util.ToBytes(typ)) {
			// 		return "types.Votes", nil
			// 	}
			case *types.Price:
				if bytes.Equal(v, util.ToBytes(typ)) {
					return "types.Price", nil
				}
			case *types.PriceFeed:
				if bytes.Equal(v, util.ToBytes(typ)) {
					return "types.PriceFeed", nil
				}
			case *types.AssetFeed:
				if bytes.Equal(v, util.ToBytes(typ)) {
					return "types.AssetFeed", nil
				}
			}
			// if val1, ok := value.(map[string]interface{}); ok {
			// 	val2 := util.ToMap(typ)

			// 	util.Dump("compare-1", val1)
			// 	util.Dump("compare-2", val2)

			// 	if reflect.DeepEqual(val1, val2) {
			// 		util.Dump("solved", typ)
			// 		return "lala", nil
			// 	}
			// }
			// name := reflect.ValueOf(typ).Type().Name()
			// util.Dump("solved", typ)
			// util.Dump("typeName", name)
			//o1 := util.ToMap(value)
		}
	}

	return suggestedType, nil
}
