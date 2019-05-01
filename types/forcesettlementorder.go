package types

//go:generate ffjson $GOFILE

type ForceSettlementOrders []ForceSettlementOrder

type ForceSettlementOrder struct {
	ID             ForceSettlementID `json:"id"`
	Owner          AccountID         `json:"owner"`
	SettlementDate Time              `json:"settlement_date"`
	Balance        AssetAmount       `json:"balance"`
}
