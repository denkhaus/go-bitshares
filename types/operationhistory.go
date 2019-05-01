package types

//go:generate ffjson $GOFILE

type OperationHistories []OperationHistory
type OperationRelativeHistories []OperationRelativeHistory

type OperationHistory struct {
	ID         OperationHistoryID `json:"id"`
	BlockNum   UInt32             `json:"block_num"`
	TrxInBlock UInt16             `json:"trx_in_block"`
	OpInTrx    UInt16             `json:"op_in_trx"`
	VirtualOp  UInt16             `json:"virtual_op"`
	Operation  OperationEnvelope  `json:"op"`
	Result     OperationResult    `json:"result"`
}

type OperationRelativeHistory struct {
	Memo        Buffer           `json:"memo"`
	Description String           `json:"description"`
	Op          OperationHistory `json:"op"`
}
