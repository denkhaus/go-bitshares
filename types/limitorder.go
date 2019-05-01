package types

//go:generate ffjson $GOFILE

type LimitOrders []LimitOrder

type LimitOrder struct {
	ID          LimitOrderID `json:"id"`
	Seller      AccountID    `json:"seller"`
	Expiration  Time         `json:"expiration"`
	ForSale     UInt64       `json:"for_sale"`
	DeferredFee UInt64       `json:"deferred_fee"`
	SellPrice   Price        `json:"sell_price"`
}
