package types

type OperationFee struct {
	Fee *AssetAmount `json:"fee,omitempty"`
}

func (p OperationFee) GetFee() AssetAmount {
	return *p.Fee
}

func (p *OperationFee) SetFee(fee AssetAmount) {
	p.Fee = &fee
}
