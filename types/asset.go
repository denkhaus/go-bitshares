package types

//go:generate ffjson   $GOFILE

type Assets []Asset

type Asset struct {
	ID                 GrapheneID   `json:"id"`
	Symbol             string       `json:"symbol"`
	Precision          int          `json:"precision"`
	Issuer             GrapheneID   `json:"issuer"`
	DynamicAssetDataID GrapheneID   `json:"dynamic_asset_data_id"`
	BitassetDataID     GrapheneID   `json:"bitasset_data_id"`
	Options            AssetOptions `json:"options"`
}

//NewAsset creates a new Asset object
func NewAsset(id GrapheneID) *Asset {
	ass := Asset{
		ID: id,
	}

	return &ass
}
