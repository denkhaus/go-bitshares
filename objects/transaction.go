package objects

import (
	"bytes"

	"encoding/hex"

	"time"

	"github.com/denkhaus/bitshares/crypto"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

//go:generate ffjson $GOFILE

type Signature string
type Signatures []Signature

const (
	TransactionExpirationTime = 30
)

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

//Sign signes a Transaction with the given private keys
func (p *Transaction) Sign(privKeys [][]byte, props *DynamicGlobalProperties, chainID string) error {

	//set Block data
	prefix, err := props.RefBlockPrefix()
	if err != nil {
		return errors.Annotate(err, "RefBlockPrefix")
	}

	p.RefBlockPrefix = prefix
	p.RefBlockNum = props.RefBlockNum()
	p.Expiration = props.Time.Add(30 * time.Second)

	var buf bytes.Buffer
	enc := util.NewTypeEncoder(&buf)

	rawChainID, err := hex.DecodeString(chainID)
	if err != nil {
		return errors.Annotatef(err, "decode chainID: %v", chainID)
	}

	if err := enc.Encode(rawChainID); err != nil {
		return errors.Annotate(err, "encode chainID")
	}

	if err := enc.Encode(p); err != nil {
		return errors.Annotate(err, "encode transaction")
	}

	data := buf.Bytes()
	p.Signatures = make([]Signature, len(privKeys))

	for idx, key := range privKeys {
		sig, err := crypto.Sign(key, data)
		if err != nil {
			return errors.Annotate(err, "Sign")
		}
		p.Signatures[idx] = Signature(hex.EncodeToString(sig))
	}

	return nil
}

//NewTransaction creates a new Transaction
func NewTransaction() *Transaction {
	tx := Transaction{
		Extensions: []Extension{},
	}
	return &tx
}
