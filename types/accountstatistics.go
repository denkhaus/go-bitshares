package types

//go:generate ffjson $GOFILE

type AccountStatistics struct {
	ID                GrapheneID `json:"id"`
	MostRecentOp      GrapheneID `json:"most_recent_op"`
	Owner             GrapheneID `json:"owner"`
	LifetimeFeesPaid  UInt64     `json:"lifetime_fees_paid"`
	PendingFees       UInt64     `json:"pending_fees"`
	PendingVestedFees UInt64     `json:"pending_vested_fees"`
	RemovedOps        UInt64     `json:"removed_ops"`
	TotalOps          UInt64     `json:"total_ops"`
	TotalCoreInOrders UInt64     `json:"total_core_in_orders"`
}
