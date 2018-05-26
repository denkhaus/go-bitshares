package types

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/juju/errors"
)

var (
	ErrInvalidCurve               = fmt.Errorf("invalid elliptic curve")
	ErrSharedKeyTooBig            = fmt.Errorf("shared key params are too big")
	ErrSharedKeyIsPointAtInfinity = fmt.Errorf("shared key is point at infinity")
)

type PrivateKey struct {
	priv *btcec.PrivateKey
	pub  *PublicKey
}

func NewPrivateKeyFromWif(wifPrivateKey string) (*PrivateKey, error) {
	w, err := btcutil.DecodeWIF(wifPrivateKey)
	if err != nil {
		return nil, errors.Annotate(err, "DecodeWIF")
	}

	priv := w.PrivKey
	k := PrivateKey{
		priv: priv,
		pub:  NewPublicKey(priv.PubKey()),
	}

	return &k, nil
}

func (p PrivateKey) PublicKey() *PublicKey {
	return p.pub
}

func (p PrivateKey) Bytes() []byte {
	return p.priv.Serialize()
}

func (p PrivateKey) SharedSecret(pub *PublicKey, skLen, macLen int) (sk []byte, err error) {
	pk := pub.ToECDSA()
	if p.priv.PublicKey.Curve != pk.Curve {
		return nil, ErrInvalidCurve
	}

	if skLen+macLen > pub.MaxSharedKeyLength() {
		return nil, ErrSharedKeyTooBig
	}

	x, _ := pk.Curve.ScalarMult(pk.X, pk.Y, p.priv.D.Bytes())
	if x == nil {
		return nil, ErrSharedKeyIsPointAtInfinity
	}

	sk = make([]byte, skLen+macLen)
	skBytes := x.Bytes()
	copy(sk[len(sk)-len(skBytes):], skBytes)
	return sk, nil
}
