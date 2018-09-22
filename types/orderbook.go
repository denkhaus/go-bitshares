package types

//go:generate ffjson $GOFILE

type OrderBook struct {
	Base  GrapheneID `json:"base"`
	Asks  []Order    `json:"asks"`
	Quote GrapheneID `json:"quote"`
	Bids  []Order    `json:"bids"`
}

type Order struct {
	Base  Float64 `json:"base"`
	Quote Float64 `json:"quote"`
	Price Float64 `json:"price"`
}
