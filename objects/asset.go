package objects

type Asset struct {
	GrapheneID
	Symbol             string       `json:"symbol"`
	Precision          int          `json:"precision"`
	Issuer             string       `json:"issuer"`
	Description        string       `json:"description"`
	DynamicAssetDataID string       `json:"dynamic_asset_data_id"`
	Options            AssetOptions `json:"options"`
	BitassetDataID     string       `json:"bitasset_data_id"`
	typ                AssetType
}

func NewAsset(id ObjectID) *Asset {
	asset := Asset{}
	asset.ID = id
	return &asset
}
