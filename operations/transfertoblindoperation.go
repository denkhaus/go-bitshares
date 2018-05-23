package operations

//skip this operation for now
//go:generate ffjson $GOFILE

// import (
// 	"encoding/hex"

// 	"github.com/denkhaus/bitshares/types"
// 	"github.com/denkhaus/bitshares/util"
// 	"github.com/juju/errors"
// )

// func init() {
// 	op := &TransferToBlindOperation{}
// 	types.OperationMap[op.Type()] = op
// }

// type TransferToBlindOperationOutputs []TransferToBlindOperationOutput

// func (p TransferToBlindOperationOutputs) Marshal(enc *util.TypeEncoder) error {
// 	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
// 		return errors.Annotate(err, "encode length")
// 	}

// 	for _, o := range p {
// 		if err := enc.Encode(o); err != nil {
// 			return errors.Annotate(err, "encode Output")
// 		}
// 	}

// 	return nil
// }

// type TransferToBlindOperationOutput struct {
// 	Commitment string          `json:"commitment"`
// 	Owner      types.Authority `json:"owner"`
// 	RangeProof string          `json:"range_proof"`
// }

// //TODO: validate order
// func (p TransferToBlindOperationOutput) Marshal(enc *util.TypeEncoder) error {

// 	comm, err := hex.DecodeString(p.Commitment)
// 	if err != nil {
// 		return errors.Annotate(err, "DecodeString")
// 	}

// 	if err := enc.Encode(comm); err != nil {
// 		return errors.Annotate(err, "encode Commitment")
// 	}

// 	proof, err := hex.DecodeString(p.RangeProof)
// 	if err != nil {
// 		return errors.Annotate(err, "DecodeString")
// 	}

// 	if err := enc.Encode(proof); err != nil {
// 		return errors.Annotate(err, "encode RangeProof")
// 	}

// 	// if err := enc.Encode(p.Owner); err != nil {
// 	// 	return errors.Annotate(err, "encode Owner")
// 	// }

// 	return nil
// }

// type TransferToBlindOperation struct {
// 	Amount         types.AssetAmount               `json:"amount"`
// 	BlindingFactor string                          `json:"blinding_factor"`
// 	Fee            types.AssetAmount               `json:"fee"`
// 	From           types.GrapheneID                `json:"from"`
// 	Outputs        TransferToBlindOperationOutputs `json:"outputs"`
// }

// func (p *TransferToBlindOperation) ApplyFee(fee types.AssetAmount) {
// 	p.Fee = fee
// }

// func (p TransferToBlindOperation) Type() types.OperationType {
// 	return types.OperationTypeTransferToBlind
// }

// //TODO: validate order
// func (p TransferToBlindOperation) Marshal(enc *util.TypeEncoder) error {
// 	if err := enc.Encode(int8(p.Type())); err != nil {
// 		return errors.Annotate(err, "encode OperationType")
// 	}

// 	if err := enc.Encode(p.Fee); err != nil {
// 		return errors.Annotate(err, "encode Fee")
// 	}

// 	if err := enc.Encode(p.Amount); err != nil {
// 		return errors.Annotate(err, "encode Amount")
// 	}

// 	if err := enc.Encode(p.From); err != nil {
// 		return errors.Annotate(err, "encode From")
// 	}

// 	factor, err := hex.DecodeString(p.BlindingFactor)
// 	if err != nil {
// 		return errors.Annotate(err, "DecodeString")
// 	}

// 	if err := enc.Encode(factor); err != nil {
// 		return errors.Annotate(err, "encode BlindingFactor")
// 	}

// 	if err := enc.Encode(p.Outputs); err != nil {
// 		return errors.Annotate(err, "encode Outputs")
// 	}

// 	return nil
// }

// //NewTransferToBlindOperation creates a new TransferToBlindOperation
// func NewTransferToBlindOperation() *TransferToBlindOperation {
// 	tx := TransferToBlindOperation{}
// 	return &tx
// }
