package objects

type AccountOptions struct {
	MemoKey       PublicKey  `json:"memo_key"`
	VotingAccount GrapheneID `json:"voting_account"`
	NumWitness    int        `json:"num_witness"`
	NumComittee   int        `json:"num_comittee"`
	Votes         Votes      `json:"votes"`
	Extensions    Extensions `json:"extensions"`
}
