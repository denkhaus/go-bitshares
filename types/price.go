package types

//go:generate ffjson   $GOFILE

type Price struct {
	Base  AssetAmount `json:"base"`
	Quote AssetAmount `json:"quote"`
}

func (p Price) Rate(precBase, precQuote float64) Rate {
	return Rate(p.Base.Rate(precBase) / p.Quote.Rate(precQuote))
}

func (p Price) Valid() bool {
	return p.Base.Valid() && p.Quote.Valid()
}
