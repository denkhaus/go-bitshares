package types

import (
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"github.com/denkhaus/bitshares/config"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type PublicKey struct {
	key string
}

func (p PublicKey) String() string {
	return p.key
}

func (p *PublicKey) UnmarshalJSON(s []byte) error {
	str := string(s)

	q, err := util.SafeUnquote(str)
	if err != nil {
		return errors.Annotate(err, "SafeUnquote")
	}

	cnf := config.CurrentConfig()
	prefix := cnf.Prefix()

	if !strings.HasPrefix(q, prefix) {
		return ErrInvalidPublicKeyForThisChain
	}

	p.key = q
	return nil
}

func (p PublicKey) Marshal(enc *util.TypeEncoder) error {
	return enc.Encode(p.Bytes())
}

func (p PublicKey) Bytes() []byte {
	if len(p.key) == 0 {
		return EmptyBuffer
	}

	cnf := config.CurrentConfig()
	prefix := cnf.Prefix()

	key := p.key
	if strings.IndexAny(key, prefix) == 0 {
		key = key[len(prefix):]
	}

	b := base58.Decode(key)
	return b[:btcec.PubKeyBytesLenCompressed]
}

func NewPublicKey(key string) *PublicKey {
	k := PublicKey{key: key}
	return &k
}
