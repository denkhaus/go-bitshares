package types

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"github.com/denkhaus/bitshares/config"
	"github.com/denkhaus/bitshares/util"
	sort "github.com/emirpasic/gods/utils"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type PublicKeys []PublicKey

func (p PublicKeys) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	for _, pub := range p {
		if err := enc.Encode(pub); err != nil {
			return errors.Annotate(err, "encode PublicKey")
		}
	}

	return nil
}

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

	pub, err := NewPublicKeyFromString(key)
	if err != nil {
		return errors.Annotate(err, "NewPublicKeyFromString")
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

func (p *PublicKey) ToAddress() (*Address, error) {
	return NewAddress(p)
}

func (p PublicKey) Bytes() []byte {
	return p.key.SerializeCompressed()
}

func (p PublicKey) Equal(pub *PublicKey) bool {
	return p.key.IsEqual(pub.key)
}

func (p PublicKey) ToECDSA() *ecdsa.PublicKey {
	return p.key.ToECDSA()
}

// MaxSharedKeyLength returns the maximum length of the shared key the
// public key can produce.
func (p PublicKey) MaxSharedKeyLength() int {
	return (p.key.ToECDSA().Curve.Params().BitSize + 7) / 8
}

//NewPublicKey creates a new PublicKey from string
//e.g.("BTS6K35Bajw29N4fjP4XADHtJ7bEj2xHJ8CoY2P2s1igXTB5oMBhR")
func NewPublicKeyFromString(key string) (*PublicKey, error) {
	cnf := config.CurrentConfig()
	if cnf == nil {
		return nil, ErrCurrentChainConfigIsNotSet
	}

	prefixChain := cnf.Prefix()
	prefix := key[:len(prefixChain)]

	if prefix != prefixChain {
		return nil, ErrPublicKeyChainPrefixMismatch
	}

	b58 := base58.Decode(key[len(prefixChain):])
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

	k := PublicKey{
		key:      pub,
		prefix:   prefix,
		checksum: chk1,
	}

	return &k, nil
}

func NewPublicKey(pub *btcec.PublicKey) (*PublicKey, error) {
	buf := pub.SerializeCompressed()
	cnf := config.CurrentConfig()
	if cnf == nil {
		return nil, ErrCurrentChainConfigIsNotSet
	}

	chk, err := util.Ripemd160Checksum(buf)
	if err != nil {
		return nil, errors.Annotate(err, "Ripemd160Checksum")
	}

	k := PublicKey{
		key:      pub,
		prefix:   cnf.Prefix(),
		checksum: chk,
	}

	return &k, nil
}

func publicKeyComparator(key1, key2 *PublicKey) (int, error) {
	addr1, err := key1.ToAddress()
	if err != nil {
		return 0, errors.Annotate(err, "ToAddress 1")
	}

	addr2, err := key2.ToAddress()
	if err != nil {
		return 0, errors.Annotate(err, "ToAddress 2")
	}

	return sort.StringComparator(addr1.String(), addr2.String()), nil
}
