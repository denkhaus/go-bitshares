package types

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"

	"github.com/pquerna/ffjson/ffjson"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type PublicKeys []PublicKey

type PublicKey struct {
	key      *btcec.PublicKey
	prefix   string
	checksum []byte
}

func (p PublicKey) String() string {
	b := append(p.Bytes(), p.checksum...)
	return fmt.Sprintf("%s%s", p.prefix, base58.Encode(b))
}

func (p *PublicKey) UnmarshalJSON(data []byte) error {
	var key string

	if err := ffjson.Unmarshal(data, &key); err != nil {
		return errors.Annotate(err, "Unmarshal")
	}

	pub, err := NewPublicKey(key)
	if err != nil {
		return errors.Annotate(err, "NewPublicKey")
	}

	p.key = pub.key
	p.prefix = pub.prefix
	p.checksum = pub.checksum
	return nil
}

func (p PublicKey) MarshalJSON() ([]byte, error) {
	return ffjson.Marshal(p.String())
}

func (p PublicKey) Marshal(enc *util.TypeEncoder) error {
	return enc.Encode(p.Bytes())
}

func (p PublicKey) Bytes() []byte {
	return p.key.SerializeCompressed()
}

func (p PublicKey) ToECDSA() *ecdsa.PublicKey {
	return p.key.ToECDSA()
}

//NewPublicKey creates a new PublicKey from string
//e.g.("BTS6K35Bajw29N4fjP4XADHtJ7bEj2xHJ8CoY2P2s1igXTB5oMBhR")
func NewPublicKey(key string) (*PublicKey, error) {
	prefix := key[:3]

	b58 := base58.Decode(key[3:])
	if len(b58) < 5 {
		return nil, ErrInvalidPublicKey
	}

	chk1 := b58[len(b58)-4:]

	keyBytes := b58[:len(b58)-4]
	chk2, err := util.Ripemd160Checksum(keyBytes)
	if err != nil {
		return nil, errors.Annotate(err, "Ripemd160Checksum")
	}

	if !bytes.Equal(chk1, chk2) {
		return nil, ErrInvalidPublicKey
	}

	pub, err := btcec.ParsePubKey(keyBytes, btcec.S256())
	if err != nil {
		return nil, errors.Annotate(err, "ParsePubKey")
	}

	k := &PublicKey{
		key:      pub,
		prefix:   prefix,
		checksum: chk1,
	}

	return k, nil
}
