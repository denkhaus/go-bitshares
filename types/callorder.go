package types

//go:generate ffjson   $GOFILE

type CallOrder struct {
	ID         GrapheneID `json:"id"`
	Borrower   GrapheneID `json:"borrower"`
	Collateral UInt64     `json:"collateral"`
	Debt       UInt64     `json:"debt"`
	CallPrice  Price      `json:"call_price"`
}
