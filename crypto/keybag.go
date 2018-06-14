package crypto

import (
	"bufio"
	"os"
	"strings"

	"github.com/denkhaus/bitshares/types"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

// KeyBag is a PrivateKey collection for signing and verifying purposes.
type KeyBag struct {
	keys []*types.PrivateKey
}

func NewKeyBag() *KeyBag {
	bag := KeyBag{
		keys: make([]*types.PrivateKey, 0),
	}

	return &bag
}

func (p KeyBag) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p.keys))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	for _, k := range p.keys {
		if err := enc.Encode(k.Bytes()); err != nil {
			return errors.Annotate(err, "encode Key")
		}
	}

	return nil
}

func (b *KeyBag) Add(wifKey string) error {
	privKey, err := types.NewPrivateKeyFromWif(wifKey)
	if err != nil {
		return errors.Annotate(err, "NewPrivateKeyFromWif")
	}

	b.keys = append(b.keys, privKey)
	return nil
}

func (b *KeyBag) ImportFromFile(path string) error {
	inFile, err := os.Open(path)
	if err != nil {
		return errors.Errorf("import keys from file [%s], %s", path, err)
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		key := strings.TrimSpace(strings.Split(scanner.Text(), " ")[0])

		if strings.Contains(key, "/") || strings.Contains(key, "#") || strings.Contains(key, ";") {
			return errors.Errorf("lines should consist of a private key on each line, with an optional whitespace and comment")
		}

		if err := b.Add(key); err != nil {
			return err
		}
	}

	return nil
}

func (b KeyBag) Publics() (out types.PublicKeys) {
	for _, k := range b.keys {
		pub := k.PublicKey()
		out = append(out, *pub)
	}
	return
}

func (b KeyBag) PublicPresent(pub *types.PublicKey) bool {
	for _, k := range b.keys {
		if k.PublicKey().Equal(pub) {
			return true
		}
	}
	return false
}

func (b KeyBag) Privates() (out types.PrivateKeys) {
	for _, k := range b.keys {
		priv := k
		out = append(out, *priv)
	}

	return
}
func (b KeyBag) PrivatesByPublics(pubKeys types.PublicKeys) (out types.PrivateKeys) {
	for _, pub := range pubKeys {
		for _, k := range b.keys {
			if pub.Equal(k.PublicKey()) {
				priv := k
				out = append(out, *priv)
			}
		}
	}

	return
}
