package types

//go:generate ffjson $GOFILE

type CallOrders []CallOrder
type CallOrder struct {
	ID         GrapheneID `json:"id"`
	Borrower   GrapheneID `json:"borrower"`
	Collateral Int64      `json:"collateral"`
	Debt       Int64      `json:"debt"`
	CallPrice  Price      `json:"call_price"`
}
