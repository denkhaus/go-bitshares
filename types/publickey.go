package types

import (
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"github.com/denkhaus/bitshares/config"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type PublicKey string

func (p *PublicKey) UnmarshalJSON(s []byte) error {
	str := string(s)

	q, err := util.SafeUnquote(str)
	if err != nil {
		return errors.Annotate(err, "SafeUnquote")
	}

	*p = PublicKey(q)
	return nil

	//return errors.Errorf("unmarshal PublicKey from %s", str)
}

func (p PublicKey) Marshal(enc *util.TypeEncoder) error {
	return enc.Encode(p.Bytes())
}

func (p PublicKey) Bytes() []byte {
	if len(p) == 0 {
		return EmptyBuffer
	}

	key := string(p)
	cnf := config.CurrentConfig()
	prefix := cnf.Prefix()

	if strings.IndexAny(key, prefix) == 0 {
		key = key[len(prefix):]
	}

	b := base58.Decode(key)
	return b[:btcec.PubKeyBytesLenCompressed]
}

func NewPublicKey(key string) *PublicKey {
	k := PublicKey(key)
	return &k
}
