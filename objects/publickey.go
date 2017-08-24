package objects

import (
	"encoding/json"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type PublicKey struct {
	key string
}

func (p *PublicKey) UnmarshalJSON(s []byte) error {
	str := string(s)

	if len(str) > 0 && str != "null" {
		q, err := util.SafeUnquote(str)
		if err != nil {
			return errors.Annotate(err, "SafeUnquote")
		}

		p.key = q
		return nil
	}

	return errors.Errorf("unmarshal PublicKey from %s", str)
}

//implements TypeMarshaller interface
func (p PublicKey) Marshal(enc *util.TypeEncoder) error {
	return nil
}

func (p PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.key)
}

func NewPublicKey(key string) *PublicKey {
	k := PublicKey{
		key: key,
	}

	return &k
}
