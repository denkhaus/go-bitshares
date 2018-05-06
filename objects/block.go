package objects

type Block struct {
	Witness               GrapheneID   `json:"witness"`
	TransactionMerkleRoot string       `json:"transaction_merkle_root"`
	WitnessSignature      string       `json:"witness_signature"`
	Previous              string       `json:"previous"`
	Extensions            Extensions   `json:"extensions"`
	TimeStamp             Time         `json:"timestamp"`
	Transactions          Transactions `json:"transactions"`
}
