package crypto

import (
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/denkhaus/bitshares/config"
	"github.com/denkhaus/bitshares/types"

	"github.com/juju/errors"
)

//TransactionSigner can sign and verify a transaction.
type TransactionSigner struct {
	*types.SignedTransaction
}

//NewTransactionSigner creates an New TransactionSigner. Invalid expiration time will be adjusted.
func NewTransactionSigner(tx *types.SignedTransaction) *TransactionSigner {
	tm := time.Now().UTC()
	if tx.Expiration.IsZero() || tx.Expiration.Before(tm) {
		tx.Expiration.Set(30 * time.Second)
	}

	return &TransactionSigner{
		SignedTransaction: tx,
	}
}

//Sign signs the underlying transaction
func (tx *TransactionSigner) Sign(privKeys types.PrivateKeys, chain *config.ChainConfig) error {
	for _, prv := range privKeys {
		ecdsaKey := prv.ToECDSA()
		if ecdsaKey.Curve != btcec.S256() {
			return types.ErrInvalidPrivateKeyCurve
		}

		for {
			digest, err := tx.Digest(chain)
			if err != nil {
				return errors.Annotate(err, "Digest")
			}

			sig, err := prv.SignCompact(digest)
			if err != nil {
				return errors.Annotate(err, "SignCompact")
			}

			if !isCanonical(sig) {
				//make canonical by adjusting expiration time
				tx.AdjustExpiration(time.Second)
			} else {
				tx.Signatures = append(tx.Signatures, types.Buffer(sig))
				break
			}
		}
	}

	return nil
}

//Verify verifies the underlying transaction against a given KeyBag
func (tx *TransactionSigner) Verify(keyBag *KeyBag, chain *config.ChainConfig) (bool, error) {
	dig, err := tx.Digest(chain)
	if err != nil {
		return false, errors.Annotate(err, "Digest")
	}

	pubKeysFound := make([]*types.PublicKey, 0, len(tx.Signatures))
	for _, signature := range tx.Signatures {
		sig := signature.Bytes()

		p, _, err := btcec.RecoverCompact(btcec.S256(), sig, dig)
		if err != nil {
			return false, errors.Annotate(err, "RecoverCompact")
		}

		pub, err := types.NewPublicKey(p)
		if err != nil {
			return false, errors.Annotate(err, "NewPublicKey")
		}

		pubKeysFound = append(pubKeysFound, pub)
	}

	for _, pub := range pubKeysFound {
		if !keyBag.Present(pub) {
			return false, nil
		}
	}

	return true, nil
}

func isCanonical(sig []byte) bool {
	d := sig
	t1 := (d[1] & 0x80) == 0
	t2 := !(d[1] == 0 && ((d[2] & 0x80) == 0))
	t3 := (d[33] & 0x80) == 0
	t4 := !(d[33] == 0 && ((d[34] & 0x80) == 0))
	return t1 && t2 && t3 && t4
}
