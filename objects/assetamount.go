package objects

import (
	json "encoding/json"
	"strconv"

	"github.com/juju/errors"
	jlexer "github.com/mailru/easyjson/jlexer"
)

type AssetAmount struct {
	Asset  GrapheneID `json:"asset_id"`
	Amount uint64     `json:"amount"`
}

//Add adds two asset amounts. They must refer to the same Asset type.
//other: The other AssetAmount to add to this.
//return: The same instance of the AssetAmount class with the combined amount.
func (p *AssetAmount) Add(other AssetAmount) *AssetAmount {
	if p.Asset.Id() != other.Asset.Id() {
		panic("Cannot add two AssetAmount instances that refer to different assets")
	}

	p.Amount += other.Amount
	return p
}

//Subtract subtracts another instance of AssetAmount from this one. This method will always
//return absolute values.
//other: The other asset amount to subtract from this.
//return: The same instance of the AssetAmount class with the combined amount.
func (p *AssetAmount) Subtract(other AssetAmount) *AssetAmount {
	if p.Asset.Id() != other.Asset.Id() {
		panic("Cannot subtract two AssetAmount instances that refer to different assets")
	}

	if p.Amount > other.Amount {
		p.Amount -= other.Amount
	} else {
		p.Amount = other.Amount - p.Amount
	}

	return p
}

func (p *AssetAmount) UnmarshalEasyJSON(j *jlexer.Lexer) {
	var res map[string]interface{}
	if err := json.Unmarshal(j.Data, &res); err != nil {
		panic(errors.Annotate(err, "unmarshal AssetAmount").Error())
	}

	p.Asset = *NewGrapheneID(ObjectID(res["asset_id"].(string)))

	if am, ok := res["amount"].(string); ok {
		amount, err := strconv.ParseUint(am, 10, 64)
		if err != nil {
			panic(errors.Annotate(err, "parse AssetAmount [amount]").Error())
		}

		p.Amount = amount
	} else {
		p.Amount = uint64(res["amount"].(float64))
	}
}
