package types

//go:generate ffjson $GOFILE

type Block struct {
	Witness               GrapheneID         `json:"witness"`
	TransactionMerkleRoot Buffer             `json:"transaction_merkle_root"`
	WitnessSignature      Buffer             `json:"witness_signature"`
	Previous              Buffer             `json:"previous"`
	BlockID               Buffer             `json:"block_id"`
	TimeStamp             Time               `json:"timestamp"`
	SigningKey            PublicKey          `json:"signing_key"`
	Transactions          SignedTransactions `json:"transactions"`
	TransactionIDs        Buffers            `json:"transaction_ids"`
	Extensions            Extensions         `json:"extensions"`
}
