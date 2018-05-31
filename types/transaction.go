package types

import (
	"bytes"
	"encoding/hex"
	"time"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

//go:generate ffjson $GOFILE

const (
	TxExpirationDefault = 30 * time.Second
)

//TODO: implement
// @property
// def id(self):
// 	""" The transaction id of this transaction
// 	"""
// 	# Store signatures temporarily since they are not part of
// 	# transaction id
// 	sigs = self.data["signatures"]
// 	self.data.pop("signatures", None)

// 	# Generage Hash of the seriliazed version
// 	h = hashlib.sha256(bytes(self)).digest()

// 	# recover signatures
// 	self.data["signatures"] = sigs

// 	# Return properly truncated tx hash
// return hexlify(h[:20]).decode("ascii")

type Transactions []Transaction

type Transaction struct {
	RefBlockNum    UInt16     `json:"ref_block_num"`
	RefBlockPrefix UInt32     `json:"ref_block_prefix"`
	Expiration     Time       `json:"expiration"`
	Operations     Operations `json:"operations"`
	Signatures     Signatures `json:"signatures"`
	Extensions     Extensions `json:"extensions"`
}

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

	if err := enc.Encode(p.Signatures); err != nil {
		return errors.Annotate(err, "encode Signatures")
	}

	return nil
}

//AdjustExpiration extends expiration by given duration.
func (p *Transaction) AdjustExpiration(dur time.Duration) {
	p.Expiration = p.Expiration.Add(dur)
}
func (p Transaction) Bytes() []byte {
	var b bytes.Buffer

	enc := util.NewTypeEncoder(&b)
	if err := enc.Encode(p); err != nil {
		panic(errors.Annotate(err, "Encode"))
	}

	return b.Bytes()
}

func (p Transaction) ToHex() string {
	return hex.EncodeToString(p.Bytes())
}

//NewTransactionWithBlockData creates a new Transaction and initializes
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
