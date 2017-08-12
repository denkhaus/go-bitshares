package objects

//go:generate ffjson -force-regenerate $GOFILE

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
func NewAsset(id ObjectID) *Asset {
	asset := Asset{
		ID: *NewGrapheneID(id),
	}

	return &asset
}
