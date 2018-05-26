package util

import (
	"github.com/btcsuite/btcutil"
	"golang.org/x/crypto/ripemd160"

	"github.com/juju/errors"
)

// Decode can be used to turn WIF into a raw private key (32 bytes).
func Decode(wif string) ([]byte, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, errors.Annotate(err, "DecodeWIF")
	}

	return w.PrivKey.Serialize(), nil
}

// GetPublicKey returns the public key associated with the given WIF
// in the 33-byte compressed format.
func GetPublicKey(wif string) ([]byte, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, errors.Annotate(err, "DecodeWIF")
	}

	return w.PrivKey.PubKey().SerializeCompressed(), nil
}

func Ripemd160Checksum(in []byte) ([]byte, error) {
	h := ripemd160.New()

	if _, err := h.Write(in); err != nil {
		return nil, errors.Annotate(err, "Write")
	}

	sum := h.Sum(nil)
	return sum[:4], nil
}
