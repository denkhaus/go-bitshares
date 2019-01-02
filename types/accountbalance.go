package types

//go:generate ffjson $GOFILE

type AccountBalances []AccountBalance

type AccountBalance struct {
	ID              AccountBalanceID `json:"id"`
	Owner           AccountID        `json:"owner"`
	AssetType       AssetID          `json:"asset_type"`
	Balance         UInt64           `json:"balance"`
	MaintenanceFlag bool             `json:"maintenance_flag"`
}
