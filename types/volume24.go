package types

//go:generate ffjson $GOFILE

type Volume24 struct {
	Base        GrapheneID `json:"base"`
	BaseVolume  Float64    `json:"base_volume"`
	Quote       GrapheneID `json:"quote"`
	QuoteVolume Float64    `json:"quote_volume"`
	Time        Time       `json:"time"`
}
