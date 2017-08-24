package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/objects"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

//LimitOrderCreateOperation instructs the blockchain to attempt to sell one asset for another.
//The blockchain will atempt to sell amount_to_sell.asset_id for as much min_to_receive.asset_id as possible.
//The fee will be paid by the seller’s account. Market fees will apply as specified by the issuer of both the selling asset and the receiving asset as a percentage of the amount exchanged.
//If either the selling asset or the receiving asset is white list restricted, the order will only be created if the seller is on the white list of the restricted asset type.
//Market orders are matched in the order they are included in the block chain.
type LimitOrderCreateOperation struct {
	Seller       objects.GrapheneID  `json:"seller"`
	Fee          objects.AssetAmount `json:"fee"`
	AmountToSell objects.AssetAmount `json:"amount_to_sell"`
	MinToReceive objects.AssetAmount `json:"min_to_receive"`
	Expiration   objects.Time        `json:"expiration"`
	FillOrKill   bool                `json:"fill_or_kill"`
	Extensions   objects.Extensions  `json:"extensions"`
}

//implements Operation interface
func (p *LimitOrderCreateOperation) ApplyFee(fee objects.AssetAmount) {
	p.Fee = fee
}

//implements Operation interface
func (p LimitOrderCreateOperation) Type() objects.OperationType {
	return objects.OperationTypeLimitOrderCreate
}

//implements Operation interface
func (p LimitOrderCreateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation type")
	}

	return nil
}

func NewLimitOrderCreateOperation() *LimitOrderCreateOperation {
	op := LimitOrderCreateOperation{
		Extensions: objects.Extensions{},
	}

	return &op
}
