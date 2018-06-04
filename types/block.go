package types

//go:generate ffjson $GOFILE

type Block struct {
	Witness               GrapheneID         `json:"witness"`
	TransactionMerkleRoot Buffer             `json:"transaction_merkle_root"`
	WitnessSignature      Buffer             `json:"witness_signature"`
	Previous              Buffer             `json:"previous"`
	TimeStamp             Time               `json:"timestamp"`
	Transactions          SignedTransactions `json:"transactions"`
	Extensions            Extensions         `json:"extensions"`
}
