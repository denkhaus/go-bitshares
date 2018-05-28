package util

import (
	"crypto/sha512"

	"golang.org/x/crypto/ripemd160"

	"github.com/juju/errors"
)

// Decode can be used to turn WIF into a raw private key (32 bytes).
// func Decode(wif string) ([]byte, error) {
// 	w, err := btcutil.DecodeWIF(wif)
// 	if err != nil {
// 		return nil, errors.Annotate(err, "DecodeWIF")
// 	}

// 	return w.PrivKey.Serialize(), nil
// }

// // GetPublicKey returns the public key associated with the given WIF
// // in the 33-byte compressed format.
// func GetPublicKey(wif string) ([]byte, error) {
// 	w, err := btcutil.DecodeWIF(wif)
// 	if err != nil {
// 		return nil, errors.Annotate(err, "DecodeWIF")
// 	}

// 	return w.PrivKey.PubKey().SerializeCompressed(), nil
// }

func Ripemd160(in []byte) ([]byte, error) {
	h := ripemd160.New()

	if _, err := h.Write(in); err != nil {
		return nil, errors.Annotate(err, "Write")
	}

	sum := h.Sum(nil)
	return sum, nil
}

func Ripemd160Checksum(in []byte) ([]byte, error) {
	buf, err := Ripemd160(in)
	if err != nil {
		return nil, errors.Annotate(err, "Ripemd160")
	}

	return buf[:4], nil
}
func Sha512Checksum(in []byte) ([]byte, error) {
	buf := sha512.Sum512(in)
	return buf[:4], nil
}
