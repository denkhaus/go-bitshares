package types

//go:generate ffjson $GOFILE

type Order struct {
	Base  Float64 `json:"base"`
	Quote Float64 `json:"quote"`
	Price Float64 `json:"price"`
}

type OrderBook struct {
	Base  AssetID `json:"base"`
	Asks  []Order `json:"asks"`
	Quote AssetID `json:"quote"`
	Bids  []Order `json:"bids"`
}
