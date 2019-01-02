package types

//go:generate ffjson $GOFILE

type Assets []Asset

type Asset struct {
	ID                 AssetID             `json:"id"`
	Symbol             String              `json:"symbol"`
	Precision          int                 `json:"precision"`
	Issuer             AccountID           `json:"issuer"`
	DynamicAssetDataID AssetDynamicDataID  `json:"dynamic_asset_data_id"`
	BitassetDataID     AssetBitAssetDataID `json:"bitasset_data_id"`
	Options            AssetOptions        `json:"options"`
}
