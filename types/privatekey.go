package types

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/juju/errors"
)

var (
	ErrInvalidCurve               = fmt.Errorf("invalid elliptic curve")
	ErrSharedKeyTooBig            = fmt.Errorf("shared key params are too big")
	ErrSharedKeyIsPointAtInfinity = fmt.Errorf("shared key is point at infinity")
)

type PrivateKeys []PrivateKey

type PrivateKey struct {
	priv *btcec.PrivateKey
	pub  *PublicKey
	raw  []byte
}

func NewPrivateKeyFromWif(wifPrivateKey string) (*PrivateKey, error) {
	w, err := btcutil.DecodeWIF(wifPrivateKey)
	if err != nil {
		return nil, errors.Annotate(err, "DecodeWIF")
	}

	priv := w.PrivKey
	raw := base58.Decode(wifPrivateKey)
	pub, err := NewPublicKey(priv.PubKey())
	if err != nil {
		return nil, errors.Annotate(err, "NewPublicKey")
	}

	k := PrivateKey{
		priv: priv,
		raw:  raw,
		pub:  pub,
	}

	return &k, nil
}

func (p PrivateKey) PublicKey() *PublicKey {
	return p.pub
}

func (p PrivateKey) ECPrivateKey() *btcec.PrivateKey {
	return p.priv
}

func (p PrivateKey) ToECDSA() *ecdsa.PrivateKey {
	return p.priv.ToECDSA()
}

func (p PrivateKey) Bytes() []byte {
	return p.priv.Serialize()
}

func (p PrivateKey) ToHex() string {
	return hex.EncodeToString(p.Bytes())
}

func (p PrivateKey) ToWIF() string {
	return base58.Encode(p.raw)
}

func (p PrivateKey) SignCompact(hash []byte) (sig []byte, err error) {
	sig, err = btcec.SignCompact(btcec.S256(), p.ECPrivateKey(), hash, true)
	return
}

func (p PrivateKey) SharedSecret(pub *PublicKey, skLen, macLen int) (sk []byte, err error) {
	puk := pub.ToECDSA()
	pvk := p.priv

	if pvk.PublicKey.Curve != puk.Curve {
		return nil, ErrInvalidCurve
	}

	if skLen+macLen > pub.MaxSharedKeyLength() {
		return nil, ErrSharedKeyTooBig
	}

	x, _ := puk.Curve.ScalarMult(puk.X, puk.Y, pvk.D.Bytes())
	if x == nil {
		return nil, ErrSharedKeyIsPointAtInfinity
	}

	sk = make([]byte, skLen+macLen)
	skBytes := x.Bytes()
	copy(sk[len(sk)-len(skBytes):], skBytes)
	return sk, nil
}
