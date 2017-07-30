package objects

//easyjson:json
type Asset struct {
	//TODO: GrapheneID doesn't deserialize properly when nested
	GrapheneID         `json:"id"`
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
		GrapheneID: *NewGrapheneID(id),
	}

	return &asset
}
