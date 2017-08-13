package objects

//go:generate ffjson $GOFILE

type LimitOrder struct {
	ID          GrapheneID  `json:"id"`
	Seller      GrapheneID  `json:"seller"`
	Expiration  RFC3339Time `json:"expiration"`
	ForSale     UInt64      `json:"for_sale"`
	DeferredFee UInt64      `json:"deferred_fee"`
	SellPrice   Price       `json:"sell_price"`
}
