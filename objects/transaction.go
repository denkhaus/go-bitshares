package objects

import (
	"bytes"
	"fmt"
	"time"

	"encoding/hex"

	"github.com/denkhaus/bitshares/crypto"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

//go:generate ffjson $GOFILE

type Signature string
type Signatures []Signature

const (
	TxExpirationDefault = 30 * time.Second
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
func (p *Transaction) Sign(wifKeys []string, chainID string) error {
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

	p.Signatures = make([]Signature, len(wifKeys))
	for idx, wif := range wifKeys {

		key, err := crypto.GetPrivateKey(wif)
		if err != nil {
			return errors.Annotate(err, "GetPrivateKey")
		}

		sig, err := crypto.Sign(data, key)
		if err != nil {
			return errors.Annotate(err, "Sign")
		}

		fmt.Print("canonical:", sig.IsCanonical())
		p.Signatures[idx] = Signature(sig.ToHex())
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
