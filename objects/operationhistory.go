package objects

//go:generate ffjson -force-regenerate $GOFILE

type OperationHistory struct {
	ID GrapheneID `json:"id"`
	// the block that caused this operation
	BlockNum UInt32 `json:"block_num"`
	// the transaction in the block
	TrxInBlock UInt16 `json:"trx_in_block"`
	// the operation within the transaction
	OpInTrx UInt16 `json:"op_in_trx"`
	// any virtual operations implied by operation in block
	VirtualOp UInt16 `json:"virtual_op"`
	// the operation
	Operation OperationEnvelope `json:"op"`
	// the operation result
	Result OperationResult `json:"result"`
}
