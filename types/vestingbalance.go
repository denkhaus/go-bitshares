package types

//go:generate ffjson $GOFILE

type VestingBalances []VestingBalance

type VestingBalance struct {
	ID      VestingBalanceID `json:"id"`
	Balance AssetAmount      `json:"balance"`
	Owner   AccountID        `json:"owner"`
	Policy  VestingPolicy    `json:"policy"`
}
