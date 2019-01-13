package types

//go:generate ffjson $GOFILE

type BroadcastResponse struct {
	ID       string            `json:"id"`
	BlockNum UInt64            `json:"block_num"`
	TrxNum   UInt32            `json:"trx_num"`
	Expired  bool              `json:"expired"`
	Trx      SignedTransaction `json:"trx"`
}
