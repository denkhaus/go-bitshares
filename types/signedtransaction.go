package types

//go:generate ffjson $GOFILE

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/denkhaus/bitshares/config"
	"github.com/denkhaus/bitshares/util"
	"github.com/denkhaus/logging"
	"github.com/pquerna/ffjson/ffjson"

	"github.com/juju/errors"
)

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

type SignedTransactions []SignedTransaction

type SignedTransaction struct {
	Transaction
	Signatures Signatures `json:"signatures"`
}

func (p SignedTransaction) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.Transaction); err != nil {
		return errors.Annotate(err, "encode Transaction")
	}

	if err := enc.Encode(p.Signatures); err != nil {
		return errors.Annotate(err, "encode Signatures")
	}

	return nil
}

//SerializeTrx serializes the transaction wihout signatures.
func (p SignedTransaction) SerializeTrx() ([]byte, error) {
	var b bytes.Buffer
	enc := util.NewTypeEncoder(&b)
	if err := enc.Encode(p.Transaction); err != nil {
		return nil, errors.Annotate(err, "encode Transaction")
	}

	return b.Bytes(), nil
}

//ToHex returns th hex representation of the underlying transaction + signatures.
func (p SignedTransaction) ToHex() (string, error) {
	var b bytes.Buffer
	enc := util.NewTypeEncoder(&b)
	if err := enc.Encode(p); err != nil {
		return "", errors.Annotate(err, "encode SignedTransaction")
	}

	return hex.EncodeToString(b.Bytes()), nil
}

//Digest calculates ths sha256 hash of the transaction.
func (tx SignedTransaction) Digest(chain *config.ChainConfig) ([]byte, error) {
	if chain == nil {
		return nil, ErrChainConfigIsUndefined
	}

	writer := sha256.New()
	rawChainID, err := hex.DecodeString(chain.ID)
	if err != nil {
		return nil, errors.Annotatef(err, "failed to decode chain ID: %v", chain.ID)
	}

	//	digestChainID := sha256.Sum256(rawChainID)
	//	util.Dump("digest chainID", hex.EncodeToString(digestChainID[:]))

	if _, err := writer.Write(rawChainID); err != nil {
		return nil, errors.Annotate(err, "Write [chainID]")
	}

	rawTrx, err := tx.SerializeTrx()
	if err != nil {
		return nil, errors.Annotatef(err, "Serialize")
	}

	//	digestTrx := sha256.Sum256(rawTrx)
	//	util.Dump("digest trx", hex.EncodeToString(digestTrx[:]))

	if _, err := writer.Write(rawTrx); err != nil {
		return nil, errors.Annotate(err, "Write [trx]")
	}

	digest := writer.Sum(nil)
	//	util.Dump("digest trx all", hex.EncodeToString(digest[:]))

	return digest[:], nil
}

//NewSignedTransactionWithBlockData creates a new SignedTransaction and initializes
//relevant Blockdata fields and expiration.
func NewSignedTransactionWithBlockData(props *DynamicGlobalProperties) (*SignedTransaction, error) {
	prefix, err := props.RefBlockPrefix()
	if err != nil {
		return nil, errors.Annotate(err, "RefBlockPrefix")
	}

	tx := SignedTransaction{
		Transaction: Transaction{
			Extensions:     Extensions{},
			RefBlockNum:    props.RefBlockNum(),
			Expiration:     props.Time.Add(TxExpirationDefault),
			RefBlockPrefix: prefix,
		},
		Signatures: Signatures{},
	}

	return &tx, nil
}

//NewSignedTransaction creates an new SignedTransaction
func NewSignedTransaction() *SignedTransaction {
	tm := time.Now().UTC().Add(TxExpirationDefault)
	tx := SignedTransaction{
		Transaction: Transaction{
			Extensions: Extensions{},
			Expiration: Time{tm},
		},
		Signatures: Signatures{},
	}

	return &tx
}

type SignedTransactionWithTransactionId struct {
	TransactionId     string
	SignedTransaction SignedTransaction
}

func (p SignedTransactionWithTransactionId) Marshal(enc *util.TypeEncoder) error {
	// type is marshaled by operation
	if err := enc.Encode(p); err != nil {
		return errors.Annotate(err, "Encode SignedTransactionWithTransactionId")
	}

	return nil
}

func (p SignedTransactionWithTransactionId) MarshalJSON() ([]byte, error) {
	return ffjson.Marshal([]interface{}{
		p.TransactionId,
		p.SignedTransaction,
	})
}

func (p *SignedTransactionWithTransactionId) UnmarshalJSON(data []byte) error {
	raw := make([]json.RawMessage, 2)
	if err := ffjson.Unmarshal(data, &raw); err != nil {
		return errors.Annotate(err, "Unmarshal [raw]")
	}

	if len(raw) != 2 {
		return ErrInvalidInputLength
	}

	if err := ffjson.Unmarshal(raw[0], &p.TransactionId); err != nil {
		return errors.Annotate(err, "Unmarshal [TransactionId]")
	}

	if err := ffjson.Unmarshal(raw[1], &p.SignedTransaction); err != nil {
		logging.DDumpUnmarshaled(
			fmt.Sprintf("TransactionId %s", p.TransactionId),
			raw[1],
		)
		return errors.Annotatef(err, "Unmarshal SignedTransaction %v", p.SignedTransaction)
	}

	return nil
}
