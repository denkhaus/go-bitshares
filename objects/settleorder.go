package objects

//go:generate ffjson $GOFILE

type SettleOrder struct {
	ID             GrapheneID  `json:"id"`
	Owner          GrapheneID  `json:"owner"`
	SettlementDate RFC3339Time `json:"settlement_date"`
	Balance        AssetAmount `json:"balance"`
}
