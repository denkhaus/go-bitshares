package types

import "time"

//go:generate ffjson $GOFILE

type Volume24 struct {
	Base        GrapheneObject `json:"base"`
	BaseVolume  Float64        `json:"base_volume"`
	Quote       GrapheneObject `json:"quote"`
	QuoteVolume Float64        `json:"quote_volume"`
	Time        time.Time      `json:"time"`
}
