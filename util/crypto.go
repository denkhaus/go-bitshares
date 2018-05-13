package util

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"

	"github.com/juju/errors"
)

// Decode can be used to turn WIF into a raw private key (32 bytes).
func Decode(wif string) ([]byte, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, errors.Annotate(err, "failed to decode WIF")
	}

	return w.PrivKey.Serialize(), nil
}

func GetPrivateKey(wif string) (*btcec.PrivateKey, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, errors.Annotate(err, "failed to decode WIF")
	}

	return w.PrivKey, nil
}

// GetPublicKey returns the public key associated with the given WIF
// in the 33-byte compressed format.
func GetPublicKey(wif string) ([]byte, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, errors.Annotate(err, "failed to decode WIF")
	}

	return w.PrivKey.PubKey().SerializeCompressed(), nil
}
