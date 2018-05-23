package util

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"golang.org/x/crypto/ripemd160"

	"github.com/juju/errors"
)

var (
	ErrInvalidCurve               = fmt.Errorf("invalid elliptic curve")
	ErrSharedKeyTooBig            = fmt.Errorf("shared key params are too big")
	ErrSharedKeyIsPointAtInfinity = fmt.Errorf("shared key is point at infinity")
)

// Decode can be used to turn WIF into a raw private key (32 bytes).
func Decode(wif string) ([]byte, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, errors.Annotate(err, "DecodeWIF")
	}

	return w.PrivKey.Serialize(), nil
}

func GetPrivateKey(wif string) (*btcec.PrivateKey, error) {
	w, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, errors.Annotate(err, "DecodeWIF")
	}

	return w.PrivKey, nil
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

// MaxSharedKeyLength returns the maximum length of the shared key the
// public key can produce.
func MaxSharedKeyLength(pub *ecdsa.PublicKey) int {
	return (pub.Curve.Params().BitSize + 7) / 8
}

func SharedSecret(priv *btcec.PrivateKey, pub *ecdsa.PublicKey, skLen, macLen int) (sk []byte, err error) {
	if priv.PublicKey.Curve != pub.Curve {
		return nil, ErrInvalidCurve
	}

	if skLen+macLen > MaxSharedKeyLength(pub) {
		return nil, ErrSharedKeyTooBig
	}

	x, _ := pub.Curve.ScalarMult(pub.X, pub.Y, priv.D.Bytes())
	if x == nil {
		return nil, ErrSharedKeyIsPointAtInfinity
	}

	sk = make([]byte, skLen+macLen)
	skBytes := x.Bytes()
	copy(sk[len(sk)-len(skBytes):], skBytes)
	return sk, nil
}
