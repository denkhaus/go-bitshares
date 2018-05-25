package types

//go:generate ffjson $GOFILE

type Balance struct {
	ID            GrapheneID  `json:"id"`
	Balance       AssetAmount `json:"balance"`
	LastClaimDate Time        `json:"last_claim_date"`
	Owner         Address     `json:"owner"`
}
