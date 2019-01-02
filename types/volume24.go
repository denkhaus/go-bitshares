package types

//go:generate ffjson $GOFILE

type Volume24 struct {
	Base        AssetID `json:"base"`
	BaseVolume  Float64 `json:"base_volume"`
	Quote       AssetID `json:"quote"`
	QuoteVolume Float64 `json:"quote_volume"`
	Time        Time    `json:"time"`
}
