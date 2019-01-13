package types

//go:generate ffjson $GOFILE

type BlockHeader struct {
	TransactionMerkleRoot Buffer     `json:"transaction_merkle_root"`
	Previous              Buffer     `json:"previous"`
	TimeStamp             Time       `json:"timestamp"`
	Witness               WitnessID  `json:"witness"`
	Extensions            Extensions `json:"extensions"`
}

type Block struct {
	Witness               WitnessID          `json:"witness"`
	TransactionMerkleRoot Buffer             `json:"transaction_merkle_root"`
	WitnessSignature      Buffer             `json:"witness_signature"`
	Previous              Buffer             `json:"previous"`
	BlockID               Buffer             `json:"block_id"`
	TimeStamp             Time               `json:"timestamp"`
	SigningKey            *PublicKey         `json:"signing_key,omitempty"`
	Transactions          SignedTransactions `json:"transactions"`
	TransactionIDs        Buffers            `json:"transaction_ids"`
	Extensions            Extensions         `json:"extensions"`
}
