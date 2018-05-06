package objects

import (
	"time"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

//go:generate ffjson $GOFILE

type Signature string
type Signatures []Signature

const (
	TxExpirationDefault = 30 * time.Second
)

type Transactions []Transaction

type Transaction struct {
	RefBlockNum    UInt16     `json:"ref_block_num"`
	RefBlockPrefix UInt32     `json:"ref_block_prefix"`
	Expiration     Time       `json:"expiration"`
	Operations     Operations `json:"operations"`
	Extensions     Extensions `json:"extensions"`
	Signatures     Signatures `json:"signatures"`
}

//implements TypeMarshaller interface
func (p Transaction) Marshal(enc *util.TypeEncoder) error {

	if err := enc.Encode(p.RefBlockNum); err != nil {
		return errors.Annotate(err, "encode RefBlockNum")
	}

	if err := enc.Encode(p.RefBlockPrefix); err != nil {
		return errors.Annotate(err, "encode RefBlockPrefix")
	}

	if err := enc.Encode(p.Expiration); err != nil {
		return errors.Annotate(err, "encode Expiration")
	}

	if err := enc.Encode(p.Operations); err != nil {
		return errors.Annotate(err, "encode Operations")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extension")
	}

	return nil
}

//AdjustExpiration extends expiration by given duration.
func (p *Transaction) AdjustExpiration(dur time.Duration) {
	p.Expiration = p.Expiration.Add(dur)
}

//NewTransactionWithBlockData creates a new Transaction and initialises
//relevant Blockdata fields and expiration.
func NewTransactionWithBlockData(props *DynamicGlobalProperties) (*Transaction, error) {
	prefix, err := props.RefBlockPrefix()
	if err != nil {
		return nil, errors.Annotate(err, "RefBlockPrefix")
	}

	tx := Transaction{
		Extensions:     Extensions{},
		Signatures:     Signatures{},
		RefBlockNum:    props.RefBlockNum(),
		Expiration:     props.Time.Add(TxExpirationDefault),
		RefBlockPrefix: prefix,
	}
	return &tx, nil
}

//NewTransaction creates a new Transaction
func NewTransaction() *Transaction {
	tx := Transaction{
		Extensions: Extensions{},
		Signatures: Signatures{},
	}
	return &tx
}
