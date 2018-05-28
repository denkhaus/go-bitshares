package types

import (
	"encoding/hex"
	"testing"

	"github.com/denkhaus/bitshares/config"
	"github.com/stretchr/testify/assert"
)

var privKeys = [][]string{
	{"5Hx8KiHLnc3pDLkwe2jujkTTJev72n3Qx7xtyaRNBsJDuejzh9u", "BTS5zzvbDtkbUVU1gFFsKqCE55U7JbjTp6mTh1usFv7KGgXL7HDQk"},
	{"5JyuWmopuyxFyvM9xm8fxXyujzfVnsg9cvE6z3pcib5NW1Av4rP", "BTS5yXqEBShUgcVm7Mve8Fg4RzQ2ftPpmo77aMbz884eX9aeGVvwD"},
	{"5KRZv3ZmkcE71K9KwEKG6pV6pyufkMQgCJrCu8xKLf2y7R7J8YK", "BTS5Z3vsgH6xj6HLXcsU38yo4TyoZs9AUzpfbaXbuxsAYPbutWvEP"},
	{"5JTge2oTwFqfNPhUrrm6upheByG2VXvaXBAqWdDUvK2CsygMG3Z", "BTS5shffTjVoT4J8Zrj3f2mQJw4UVKrnbx5FWYhVgov45EpBf2NYi"},
	{"5JqmjeakPoTz3ComQ7Jgg11jHxywfkJHZPhMJoBomZLrZSfRAvr", "BTS56fy8qpkLzNoguGMPgPNkkznxnx2woEg1qPq7E6gF2SeGSRyK5"},
}

var data = [][]string{
	{"5JWHY5DxTF6qN5grTtChDCYBmWHfY9zaSsw4CxEKN5eZpH9iBma", "5ad2b8df2c255d4a2996ee7d065e013e1bbb35c075ee6e5208aca44adc9a9d4c"},
	{"5KPipdRzoxrp6dDqsBfMD6oFZG356trVHV5QBGx3rABs1zzWWs8", "cf9d6121ed458f24ea456ad7ff700da39e86688988cfe5c6ed6558642cf1e32f"},
}

func Test_PrivatePublic(t *testing.T) {
	config.SetCurrentConfig(config.ChainIDBTS)

	for _, k := range privKeys {
		wif := k[0]
		pub := k[1]

		key, err := NewPrivateKeyFromWif(wif)
		if err != nil {
			assert.FailNow(t, err.Error(), "NewPrivateKeyFromWif")
		}

		assert.Equal(t, pub, key.PublicKey().String())
	}
}

func TestDecode(t *testing.T) {
	config.SetCurrentConfig(config.ChainIDBTS)

	for _, k := range data {
		wif := k[0]
		hx := k[1]

		key, err := NewPrivateKeyFromWif(wif)
		if err != nil {
			assert.FailNow(t, err.Error(), "NewPrivateKeyFromWif")
		}

		got := hex.EncodeToString(key.Bytes())
		assert.Equal(t, hx, got)
	}
}
