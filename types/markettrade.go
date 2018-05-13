package types

//go:generate ffjson   $GOFILE

type MarketTrade struct {
	DateTime Time    `json:"date"`
	Price    Float64 `json:"price"`
	Amount   Float64 `json:"amount"`
	Value    Float64 `json:"value"`
}
