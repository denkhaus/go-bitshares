package types

//go:generate ffjson $GOFILE

type LimitOrders []LimitOrder

type LimitOrder struct {
	ID          GrapheneID `json:"id"`
	Seller      GrapheneID `json:"seller"`
	Expiration  Time       `json:"expiration"`
	ForSale     UInt64     `json:"for_sale"`
	DeferredFee UInt64     `json:"deferred_fee"`
	SellPrice   Price      `json:"sell_price"`
}
