// +build !nosigning

package crypto

import (
	// Stdlib
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

type SignedTransaction struct {
	*types.Transaction
}

func NewSignedTransaction(tx *types.Transaction) *SignedTransaction {
	if tx.Expiration.IsZero() {
		exp := time.Now().Add(30 * time.Second)
		tx.Expiration = types.Time{exp}
	}

	return &SignedTransaction{tx}
}

func (tx *SignedTransaction) Serialize() ([]byte, error) {
	//remove signatures before serializing
	sigTemp := tx.Signatures
	tx.Signatures.Reset()

	defer func() {
		//and restore
		tx.Signatures = sigTemp
	}()

	return tx.Bytes(), nil
}

func (tx *SignedTransaction) Digest(chain *config.ChainConfig) ([]byte, error) {
	var msgBuffer bytes.Buffer

	// Write the chain ID.
	rawChainID, err := hex.DecodeString(chain.ID())
	if err != nil {
		return nil, errors.Annotatef(err, "failed to decode chain ID: %v", chain.ID())
	}

	if _, err := msgBuffer.Write(rawChainID); err != nil {
		return nil, errors.Annotate(err, "failed to write chain ID")
	}

	// Write the serialized transaction.
	rawTx, err := tx.Serialize()
	if err != nil {
		return nil, errors.Annotate(err, "Serialize")
	}

	if _, err := msgBuffer.Write(rawTx); err != nil {
		return nil, errors.Annotate(err, "failed to write serialized transaction")
	}

	// Compute the digest.
	digest := sha256.Sum256(msgBuffer.Bytes())
	fmt.Println(hex.EncodeToString(digest[:]))

	return digest[:], nil
}

func (tx *SignedTransaction) SignTest1(keys types.PrivateKeys, chain *config.ChainConfig) error {
	privKeys := make([][]byte, len(keys))

	for idx, k := range keys {
		privKeys[idx] = k.Bytes()
	}

	var buf bytes.Buffer
	chainid, _ := hex.DecodeString(chain.ID())
	//fmt.Println(tx.Operations[0])
	//fmt.Println(" ")
	txraw, err := tx.Serialize()
	if err != nil {
		return err
	}
	//fmt.Println(tx_raw)
	//fmt.Println(" ")
	buf.Write(chainid)
	buf.Write(txraw)
	data := buf.Bytes()
	//msg_sha := crypto.Sha256(buf.Bytes())

	var sigs types.Signatures

	for _, privb := range privKeys {
		sigBytes := tx.Sign_Single(privb, data)
		sigs = append(sigs, types.Buffer(sigBytes))
	}

	tx.Transaction.Signatures = sigs
	return nil
}

func (tx *SignedTransaction) SignTest2(privKeys types.PrivateKeys, chain *config.ChainConfig) error {
	digest, err := tx.Digest(chain)
	if err != nil {
		return errors.Annotate(err, "Digest")
	}

	for _, prv := range privKeys {
		ecdsaKey := prv.ToECDSA()
		if ecdsaKey.Curve != btcec.S256() {
			return types.ErrInvalidPrivateKeyCurve
		}

		for {
			sig, err := btcec.SignCompact(btcec.S256(), prv.ECPrivateKey(), digest, true)
			if err != nil {
				return errors.Annotate(err, "SignCompact")
			}

			if !isCanonical(sig) {
				//make canonical by adjusting expiration time
				tx.AdjustExpiration(time.Second)
				digest, err = tx.Digest(chain)
				if err != nil {
					return errors.Annotate(err, "Digest")
				}

			} else {
				tx.Signatures = append(tx.Signatures, types.Buffer(sig))
				break
			}
		}
	}

	return nil
}

func (tx *SignedTransaction) Sign(keys types.PrivateKeys, chain *config.ChainConfig) error {
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

	tx.Transaction.Signatures = si
	return nil
}

// Verify verifys the Transaction against the public keys
func (tx *SignedTransaction) Verify(pubs types.PublicKeys, chain *config.ChainConfig) (bool, error) {
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
