package objects

type AssetAmount struct {
	Asset  GrapheneID `json:"asset_id"`
	Amount UInt64     `json:"amount"`
}

//Add adds two asset amounts. They must refer to the same Asset type.
//other: The other AssetAmount to add to this.
//return: The same instance of the AssetAmount class with the combined amount.
func (p *AssetAmount) Add(other AssetAmount) *AssetAmount {
	if p.Asset.Id() != other.Asset.Id() {
		panic("Cannot add two AssetAmount instances that refer to different assets")
	}

	p.Amount += other.Amount
	return p
}

//Subtract subtracts another instance of AssetAmount from this one. This method will always
//return absolute values.
//other: The other asset amount to subtract from this.
//return: The same instance of the AssetAmount class with the combined amount.
func (p *AssetAmount) Subtract(other AssetAmount) *AssetAmount {
	if p.Asset.Id() != other.Asset.Id() {
		panic("Cannot subtract two AssetAmount instances that refer to different assets")
	}

	if p.Amount > other.Amount {
		p.Amount -= other.Amount
	} else {
		p.Amount = other.Amount - p.Amount
	}

	return p
}
