package objects

//go:generate ffjson -force-regenerate $GOFILE

type MarketTrade struct {
	DateTime Time    `json:"date"`
	Price    float64 `json:"price,string"`
	Amount   float64 `json:"amount,string"`
	Value    float64 `json:"value,string"`
}
