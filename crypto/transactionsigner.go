// +build !nosigning

package crypto

import (
	// Stdlib
	"bytes"
	"time"
	"unsafe"

	"github.com/btcsuite/btcd/btcec"
	"github.com/denkhaus/bitshares/config"
	"github.com/denkhaus/bitshares/types"

	"github.com/juju/errors"
)

// #cgo LDFLAGS: -L/usr/local/lib -lsecp256k1
// #include <stdlib.h>
// #include "signing.h"
import "C"

type TransactionSigner struct {
	*types.SignedTransaction
}

func NewTransactionSigner(tx *types.SignedTransaction) *TransactionSigner {
	if tx.Expiration.IsZero() {
		exp := time.Now().Add(30 * time.Second)
		tx.Expiration = types.Time{exp}
	}

	return &TransactionSigner{tx}
}

func (tx *TransactionSigner) SignTest2(keys types.PrivateKeys, chain *config.ChainConfig) error {
	digest, err := tx.Digest(chain)
	if err != nil {
		return errors.Annotate(err, "Digest")
	}

	var sigs types.Signatures
	for _, priv := range keys {
		sig := tx.SignSingle(priv, digest)
		sigs = append(sigs, types.Buffer(sig))
	}

	tx.Signatures = sigs
	return nil
}

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

func (tx *TransactionSigner) Verify(pubKeys types.PublicKeys, chain *config.ChainConfig) (bool, error) {
	dig, err := tx.Digest(chain)
	if err != nil {
		return false, errors.Annotate(err, "Digest")
	}

	pubKeysFound := make(types.PublicKeys, 0)
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

		pubKeysFound = append(pubKeysFound, *pub)
	}

	if len(pubKeysFound) != len(pubKeys) {
		return false, nil
	}

	for idx := range pubKeys {
		if !pubKeys[idx].Equal(&pubKeysFound[idx]) {
			return false, nil
		}
	}

	return true, nil
}

func (tx *TransactionSigner) SignTest4(keys types.PrivateKeys, chain *config.ChainConfig) error {
	privKeys := make([][]byte, len(keys))

	for idx, k := range keys {
		privKeys[idx] = k.Bytes()
	}

	digest, err := tx.Digest(chain)
	if err != nil {
		return errors.Annotate(err, "Digest")
	}

	// Sign.
	cDigest := C.CBytes(digest)
	defer C.free(cDigest)

	cKeys := make([]unsafe.Pointer, 0, len(privKeys))
	for _, key := range privKeys {
		cKeys = append(cKeys, C.CBytes(key))
	}
	defer func() {
		for _, cKey := range cKeys {
			C.free(cKey)
		}
	}()

	sigs := make([][]byte, 0, len(privKeys))
	for _, cKey := range cKeys {
		var (
			signature [64]byte
			recid     C.int
		)

		code := C.sign_transaction(
			(*C.uchar)(cDigest), (*C.uchar)(cKey), (*C.uchar)(&signature[0]), &recid)
		if code == 0 {
			return errors.New("sign_transaction returned 0")
		}

		sig := make([]byte, 65)
		sig[0] = byte(recid)
		copy(sig[1:], signature[:])

		sigs = append(sigs, sig)
	}

	// Set the signatures in the transaction.
	si := make(types.Signatures, len(sigs))
	for idx, sig := range sigs {
		si[idx] = types.Buffer(sig)
	}

	tx.Signatures = si
	return nil
}

// Verify verifys the Transaction against the public keys
func (tx *TransactionSigner) VerifyTest4(pubs types.PublicKeys, chain *config.ChainConfig) (bool, error) {
	pubKeys := make([][]byte, len(pubs))

	for idx, k := range pubs {
		pubKeys[idx] = k.Bytes()
	}

	digest, err := tx.Digest(chain)
	if err != nil {
		return false, errors.Annotate(err, "Digest")
	}

	cDigest := C.CBytes(digest)
	defer C.free(cDigest)

	// Make sure to free memory.
	cSigs := make([]unsafe.Pointer, 0, len(tx.Signatures))
	defer func() {
		for _, cSig := range cSigs {
			C.free(cSig)
		}
	}()

	// Collect verified public keys.
	pubKeysFound := make([][]byte, len(pubKeys))
	for i, signature := range tx.Signatures {
		sig := signature.Bytes()
		recoverParameter := sig[0] - 27 - 4
		sig = sig[1:]

		cSig := C.CBytes(sig)
		cSigs = append(cSigs, cSig)

		var publicKey [33]byte

		code := C.verify_recoverable_signature(
			(*C.uchar)(cDigest),
			(*C.uchar)(cSig),
			(C.int)(recoverParameter),
			(*C.uchar)(&publicKey[0]),
		)
		if code == 1 {
			pubKeysFound[i] = publicKey[:]
		}
	}

	for i := range pubKeys {
		if !bytes.Equal(pubKeysFound[i], pubKeys[i]) {
			return false, nil
		}
	}
	return true, nil
}

func isCanonical(sig []byte) bool {
	if ((sig[0] & 0x80) != 0) || (sig[0] == 0) ||
		((sig[1] & 0x80) != 0) || ((sig[32] & 0x80) != 0) ||
		(sig[32] == 0) || ((sig[33] & 0x80) != 0) {
		return false
	}

	return true
}
