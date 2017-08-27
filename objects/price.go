package objects

type Price struct {
	Base  AssetAmount `json:"base"`
	Quote AssetAmount `json:"quote"`
}

func (p Price) Rate(precBase, precQuote int) float64 {
	return p.Base.Rate(precBase) / p.Quote.Rate(precQuote)
}
