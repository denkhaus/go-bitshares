package types

//go:generate ffjson $GOFILE

type Volume24 struct {
	Base        String  `json:"base"`
	BaseVolume  Float64 `json:"base_volume"`
	Quote       String  `json:"quote"`
	QuoteVolume Float64 `json:"quote_volume"`
	Time        Time    `json:"time"`
}
