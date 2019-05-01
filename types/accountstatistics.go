package types

//go:generate ffjson $GOFILE

type AccountStatistics struct {
	ID                AccountStatisticsID         `json:"id"`
	MostRecentOp      AccountTransactionHistoryID `json:"most_recent_op"`
	Owner             AccountID                   `json:"owner"`
	LifetimeFeesPaid  UInt64                      `json:"lifetime_fees_paid"`
	PendingFees       UInt64                      `json:"pending_fees"`
	PendingVestedFees UInt64                      `json:"pending_vested_fees"`
	RemovedOps        UInt64                      `json:"removed_ops"`
	TotalOps          UInt64                      `json:"total_ops"`
	TotalCoreInOrders UInt64                      `json:"total_core_in_orders"`
}
