package types

//go:generate ffjson $GOFILE

type AccountBalances []AccountBalance

type AccountBalance struct {
	ID              GrapheneID `json:"id"`
	Owner           GrapheneID `json:"owner"`
	AssetType       GrapheneID `json:"asset_type"`
	Balance         UInt64     `json:"balance"`
	MaintenanceFlag bool       `json:"maintenance_flag"`
}
