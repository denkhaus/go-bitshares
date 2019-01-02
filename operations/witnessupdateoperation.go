package operations

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeWitnessUpdate] = func() types.Operation {
		op := &WitnessUpdateOperation{}
		return op
	}
}

type WitnessUpdateOperation struct {
	types.OperationFee
	NewSigningKey  *types.PublicKey `json:"new_signing_key,omitempty"`
	NewURL         *types.String    `json:"new_url,omitempty"`
	Witness        types.WitnessID  `json:"witness"`
	WitnessAccount types.AccountID  `json:"witness_account"`
}

func (p WitnessUpdateOperation) Type() types.OperationType {
	return types.OperationTypeWitnessUpdate
}

func (p WitnessUpdateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Witness); err != nil {
		return errors.Annotate(err, "encode new options")
	}

	if err := enc.Encode(p.WitnessAccount); err != nil {
		return errors.Annotate(err, "encode WitnessAccount")
	}

	if err := enc.Encode(p.NewURL != nil); err != nil {
		return errors.Annotate(err, "encode have NewURL")
	}

	if err := enc.Encode(p.NewURL); err != nil {
		return errors.Annotate(err, "encode NewURL")
	}

	if err := enc.Encode(p.NewSigningKey != nil); err != nil {
		return errors.Annotate(err, "encode have NewSigningKey")
	}

	if err := enc.Encode(p.NewSigningKey); err != nil {
		return errors.Annotate(err, "encode NewSigningKey")
	}

	return nil
}
