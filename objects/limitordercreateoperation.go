package objects

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

//LimitOrderCreateOperation instructs the blockchain to attempt to sell one asset for another.
//The blockchain will atempt to sell amount_to_sell.asset_id for as much min_to_receive.asset_id as possible.
//The fee will be paid by the seller’s account. Market fees will apply as specified by the issuer of both the selling asset and the receiving asset as a percentage of the amount exchanged.
//If either the selling asset or the receiving asset is white list restricted, the order will only be created if the seller is on the white list of the restricted asset type.
//Market orders are matched in the order they are included in the block chain.
type LimitOrderCreateOperation struct {
	Seller       GrapheneID  `json:"seller"`
	Fee          AssetAmount `json:"fee"`
	AmountToSell AssetAmount `json:"amount_to_sell"`
	MinToReceive AssetAmount `json:"min_to_receive"`
	Expiration   Time        `json:"expiration"`
	FillOrKill   bool        `json:"fill_or_kill"`
	Extensions   Extensions  `json:"extensions"`
}

//implements Operation interface
func (p *LimitOrderCreateOperation) ApplyFee(fee AssetAmount) {
	p.Fee = fee
}

//implements Operation interface
func (p LimitOrderCreateOperation) Type() OperationType {
	return OperationTypeLimitOrderCreate
}

//implements Operation interface
// order checked!
func (p LimitOrderCreateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode operation type")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Seller); err != nil {
		return errors.Annotate(err, "encode seller")
	}

	if err := enc.Encode(p.AmountToSell); err != nil {
		return errors.Annotate(err, "encode amount to sell")
	}

	if err := enc.Encode(p.MinToReceive); err != nil {
		return errors.Annotate(err, "encode min to receive")
	}

	if err := enc.Encode(p.Expiration); err != nil {
		return errors.Annotate(err, "encode expiration")
	}

	if err := enc.Encode(p.FillOrKill); err != nil {
		return errors.Annotate(err, "encode fill or kill")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

func NewLimitOrderCreateOperation() *LimitOrderCreateOperation {
	op := LimitOrderCreateOperation{
		Extensions: Extensions{},
	}

	return &op
}
