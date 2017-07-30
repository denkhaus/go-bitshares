package objects

// AssetAmount
type AssetAmount struct {
	amount uint64
	asset  Asset
}

// Add another AssetAmount
func (p *AssetAmount) Add(other AssetAmount) AssetAmount {
	if p.asset.Id() != other.asset.Id() {
		panic("Cannot add two AssetAmount instances that refer to different assets")
	}

	return AssetAmount{
		amount: p.amount + other.amount,
		asset:  p.asset,
	}
}
